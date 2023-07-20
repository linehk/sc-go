package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/stablecog/sc-go/database/ent"
	"github.com/stablecog/sc-go/database/repository"
	"github.com/stablecog/sc-go/log"
	"github.com/stablecog/sc-go/server/requests"
	"github.com/stablecog/sc-go/server/responses"
	"github.com/stablecog/sc-go/shared"
	"github.com/stablecog/sc-go/utils"
)

// POST generate endpoint
// Adds generate to queue, if authenticated, returns the ID of the generation
func (c *RestAPI) HandleCreateGeneration(w http.ResponseWriter, r *http.Request) {
	var user *ent.User
	if user = c.GetUserIfAuthenticated(w, r); user == nil {
		return
	}

	// Parse request body
	reqBody, _ := io.ReadAll(r.Body)
	var generateReq requests.CreateGenerationRequest
	err := json.Unmarshal(reqBody, &generateReq)
	if err != nil {
		responses.ErrUnableToParseJson(w, r)
		return
	}

	// Validation
	err = generateReq.Validate()
	if err != nil {
		responses.ErrBadRequest(w, r, err.Error())
		return
	}

	// Parse request headers
	countryCode := utils.GetCountryCode(r)
	deviceInfo := utils.GetClientDeviceInfo(r)

	// ! TODO - parallel generation toggle

	// Get model and scheduler name for cog
	modelName := shared.GetCache().GetGenerationModelNameFromID(generateReq.ModelId)
	schedulerName := shared.GetCache().GetSchedulerNameFromID(generateReq.SchedulerId)
	if modelName == "" || schedulerName == "" {
		log.Error("Error getting model or scheduler name: %s - %s", modelName, schedulerName)
		responses.ErrInternalServerError(w, r, "An unknown error has occured")
		return
	}

	// Format prompts
	generateReq.Prompt = utils.FormatPrompt(generateReq.Prompt)
	generateReq.NegativePrompt = utils.FormatPrompt(generateReq.NegativePrompt)

	// For live page update
	var livePageMsg shared.LivePageMessage
	// For keeping track of this request as it gets sent to the worker
	var requestId string
	// Cog request
	var cogReqBody requests.CogQueueRequest

	// Credits left after this operation
	var remainingCredits int

	// Wrap everything in a DB transaction
	// We do this since we want our credit deduction to be atomic with the whole process
	if err := c.Repo.WithTx(func(tx *ent.Tx) error {
		// Bind a client to the transaction
		DB := tx.Client()
		// Deduct credits from user
		start := time.Now()
		deducted, err := c.Repo.DeductCreditsFromUser(user.ID, int32(generateReq.NumOutputs), DB)
		if err != nil {
			log.Error("Error deducting credits", "err", err)
			responses.ErrInternalServerError(w, r, "Error deducting credits from user")
			return err
		} else if !deducted {
			responses.ErrInsufficientCredits(w, r)
			return responses.InsufficientCreditsErr
		}
		fmt.Printf("--- Deduct credits took took: %s\n", time.Now().Sub(start))

		remainingCredits, err = c.Repo.GetNonExpiredCreditTotalForUser(user.ID, DB)
		if err != nil {
			log.Error("Error getting remaining credits", "err", err)
			responses.ErrInternalServerError(w, r, "An unknown error has occured")
			return err
		}

		// Create generation
		start = time.Now()
		g, err := c.Repo.CreateGeneration(
			user.ID,
			string(deviceInfo.DeviceType),
			deviceInfo.DeviceOs,
			deviceInfo.DeviceBrowser,
			countryCode,
			generateReq,
			DB)
		if err != nil {
			log.Error("Error creating generation", "err", err)
			responses.ErrInternalServerError(w, r, "Error creating generation")
			return err
		}
		fmt.Printf("--- Create generation took: %s\n", time.Now().Sub(start))

		// Request Id matches generation ID
		requestId = g.ID.String()

		// For live page update
		livePageMsg = shared.LivePageMessage{
			ProcessType:      generateReq.ProcessType,
			ID:               utils.Sha256(requestId),
			CountryCode:      countryCode,
			Status:           shared.LivePageQueued,
			TargetNumOutputs: generateReq.NumOutputs,
			Width:            generateReq.Width,
			Height:           generateReq.Height,
			CreatedAt:        g.CreatedAt,
		}

		cogReqBody = requests.CogQueueRequest{
			WebhookEventsFilter: []requests.CogEventFilter{requests.CogEventFilterStart, requests.CogEventFilterStart},
			RedisPubsubKey:      shared.COG_REDIS_EVENT_CHANNEL,
			Input: requests.BaseCogRequest{
				ID:                   requestId,
				UserID:               &user.ID,
				IP:                   utils.GetIPAddress(r),
				LivePageData:         &livePageMsg,
				Prompt:               generateReq.Prompt,
				NegativePrompt:       generateReq.NegativePrompt,
				Width:                fmt.Sprint(generateReq.Width),
				Height:               fmt.Sprint(generateReq.Height),
				NumInferenceSteps:    fmt.Sprint(generateReq.InferenceSteps),
				GuidanceScale:        fmt.Sprint(generateReq.GuidanceScale),
				Model:                modelName,
				ModelId:              generateReq.ModelId,
				Scheduler:            schedulerName,
				SchedulerId:          generateReq.SchedulerId,
				Seed:                 fmt.Sprint(generateReq.Seed),
				NumOutputs:           fmt.Sprint(generateReq.NumOutputs),
				OutputImageExtension: string(shared.DEFAULT_GENERATE_OUTPUT_EXTENSION),
				OutputImageQuality:   fmt.Sprint(shared.DEFAULT_GENERATE_OUTPUT_QUALITY),
				ProcessType:          generateReq.ProcessType,
				SubmitToGallery:      generateReq.SubmitToGallery,
				MaskImageUrl:         generateReq.MaskImageUrl,
				InitImageUrl:         generateReq.InitImageUrl,
			},
		}

		start = time.Now()
		err = c.Redis.EnqueueCogRequest(r.Context(), cogReqBody)
		if err != nil {
			log.Error("Failed to write request %s to queue: %v", requestId, err)
			responses.ErrInternalServerError(w, r, "Failed to queue generate request")
			return err
		}
		fmt.Printf("--- Enqueue cog request took: %s\n", time.Now().Sub(start))
		return nil
	}); err != nil {
		log.Error("Error in transaction", "err", err)
		return
	}

	// Track the request in our internal map
	start := time.Now()
	c.Redis.SetCogRequestStreamID(r.Context(), requestId, generateReq.StreamID)
	fmt.Printf("--- Put request in map took: %s\n", time.Now().Sub(start))

	// Send live page update
	go func() {
		liveResp := repository.TaskStatusUpdateResponse{
			ForLivePage:     true,
			LivePageMessage: &livePageMsg,
		}
		respBytes, err := json.Marshal(liveResp)
		if err != nil {
			log.Error("Error marshalling sse live response", "err", err)
			return
		}
		err = c.Redis.Client.Publish(c.Redis.Client.Context(), shared.REDIS_SSE_BROADCAST_CHANNEL, respBytes).Err()
		if err != nil {
			log.Error("Failed to publish live page update", "err", err)
		}
	}()

	// Start the timeout timer
	go func() {
		// sleep
		time.Sleep(shared.REQUEST_COG_TIMEOUT)
		// this will trigger timeout if it hasnt been finished
		c.Repo.FailCogMessageDueToTimeoutIfTimedOut(requests.CogRedisMessage{
			Input:  cogReqBody.Input,
			Error:  shared.TIMEOUT_ERROR,
			Status: requests.CogFailed,
		})
	}()

	render.Status(r, http.StatusOK)
	render.JSON(w, r, &responses.TaskQueuedResponse{
		ID:               requestId,
		RemainingCredits: remainingCredits,
	})
}
