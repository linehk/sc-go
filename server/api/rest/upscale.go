package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-chi/render"
	"github.com/stablecog/sc-go/database/ent"
	"github.com/stablecog/sc-go/server/requests"
	"github.com/stablecog/sc-go/server/responses"
	"github.com/stablecog/sc-go/shared"
	"github.com/stablecog/sc-go/utils"
	"k8s.io/klog/v2"
)

func (c *RestAPI) HandleUpscale(w http.ResponseWriter, r *http.Request) {
	userID, _ := c.GetUserIDAndEmailIfAuthenticated(w, r)
	if userID == nil {
		return
	}

	// Parse request body
	reqBody, _ := io.ReadAll(r.Body)
	var upscaleReq requests.CreateUpscaleRequest
	err := json.Unmarshal(reqBody, &upscaleReq)
	if err != nil {
		responses.ErrUnableToParseJson(w, r)
		return
	}

	// Validation
	err = upscaleReq.Validate()
	if err != nil {
		responses.ErrBadRequest(w, r, err.Error())
		return
	}

	// Parse request headers
	countryCode := utils.GetCountryCode(r)
	deviceInfo := utils.GetClientDeviceInfo(r)

	// Get model name for cog
	modelName := shared.GetCache().GetUpscaleModelNameFromID(upscaleReq.ModelId)
	if modelName == "" {
		klog.Errorf("Error getting model name: %s", modelName)
		responses.ErrInternalServerError(w, r, "An unknown error has occured")
		return
	}

	// Initiate upscale
	// We need to get width/height, from our database if output otherwise from the external image
	var width int32
	var height int32

	// Image Type
	imageUrl := upscaleReq.Input
	if upscaleReq.Type == requests.UpscaleRequestTypeImage {
		width, height, err = utils.GetImageWidthHeightFromUrl(imageUrl, shared.MAX_UPSCALE_IMAGE_SIZE)
		if err != nil {
			responses.ErrBadRequest(w, r, "image_url_width_height_error")
			return
		}
	}

	// Output Type
	var outputIDStr string
	if upscaleReq.Type == requests.UpscaleRequestTypeOutput {
		outputIDStr = upscaleReq.OutputID.String()
		output, err := c.Repo.GetGenerationOutputForUser(upscaleReq.OutputID, *userID)
		if err != nil {
			if ent.IsNotFound(err) {
				responses.ErrBadRequest(w, r, "output_not_found")
				return
			}
			klog.Errorf("Error getting output: %v", err)
			responses.ErrInternalServerError(w, r, "Error getting output")
			return
		}
		if output.UpscaledImagePath != nil {
			responses.ErrBadRequest(w, r, "image_already_upscaled")
			return
		}
		imageUrl = utils.GetURLFromImagePath(output.ImagePath)

		// Get width/height of generation
		width, height, err = c.Repo.GetGenerationOutputWidthHeight(upscaleReq.OutputID)
		if err != nil {
			responses.ErrBadRequest(w, r, "Unable to retrieve width/height for upscale")
			return
		}
	}

	// For live page update
	var livePageMsg shared.LivePageMessage
	// For keeping track of this request as it gets sent to the worker
	var requestId string
	// The cog request body
	var cogReqBody requests.CogQueueRequest

	// Wrap everything in a DB transaction
	// We do this since we want our credit deduction to be atomic with the whole process
	if err := c.Repo.WithTx(func(tx *ent.Tx) error {
		// Bind transaction to client
		DB := tx.Client()

		// Charge credits
		deducted, err := c.Repo.DeductCreditsFromUser(*userID, 1, DB)
		if err != nil {
			klog.Errorf("Error deducting credits: %v", err)
			responses.ErrInternalServerError(w, r, "Error deducting credits from user")
			return err
		} else if !deducted {
			responses.ErrInsufficientCredits(w, r)
			return responses.InsufficientCreditsErr
		}

		// Create upscale
		upscale, err := c.Repo.CreateUpscale(
			*userID,
			width,
			height,
			string(deviceInfo.DeviceType),
			deviceInfo.DeviceOs,
			deviceInfo.DeviceBrowser,
			countryCode,
			upscaleReq,
			DB)
		if err != nil {
			klog.Errorf("Error creating upscale: %v", err)
			responses.ErrInternalServerError(w, r, "Error creating upscale")
			return err
		}

		// Request ID matches upscale ID
		requestId = upscale.ID.String()

		// For live page update
		livePageMsg = shared.LivePageMessage{
			ProcessType: shared.UPSCALE,
			ID:          utils.Sha256(requestId),
			CountryCode: countryCode,
			Status:      shared.LivePageQueued,
			NumOutputs:  1,
			Width:       width,
			Height:      height,
			CreatedAt:   upscale.CreatedAt,
		}

		// Send to the cog
		cogReqBody = requests.CogQueueRequest{
			WebhookEventsFilter: []requests.CogEventFilter{requests.CogEventFilterStart, requests.CogEventFilterStart},
			RedisPubsubKey:      shared.COG_REDIS_EVENT_CHANNEL,
			Input: requests.BaseCogRequest{
				ID:                   requestId,
				LivePageData:         livePageMsg,
				GenerationOutputID:   outputIDStr,
				Image:                imageUrl,
				ProcessType:          shared.UPSCALE,
				Width:                fmt.Sprint(width),
				Height:               fmt.Sprint(height),
				OutputImageExtension: string(shared.DEFAULT_UPSCALE_OUTPUT_EXTENSION),
				OutputImageQuality:   fmt.Sprint(shared.DEFAULT_UPSCALE_OUTPUT_QUALITY),
			},
		}

		err = c.Redis.EnqueueCogRequest(r.Context(), cogReqBody)
		if err != nil {
			klog.Errorf("Failed to write request %s to queue: %v", requestId, err)
			responses.ErrInternalServerError(w, r, "Failed to queue upscale request")
			return err
		}

		return nil
	}); err != nil {
		klog.Errorf("Error with transaction: %v", err)
		return
	}

	// Track the request in our internal map
	c.Redis.SetCogRequestStreamID(r.Context(), requestId, upscaleReq.StreamID)

	// Deal with live page update
	go c.Hub.BroadcastLivePageMessage(livePageMsg)

	// Start the timeout timer
	go func() {
		// sleep
		time.Sleep(shared.REQUEST_COG_TIMEOUT)
		// this will trigger timeout if it hasnt been finished
		c.Repo.FailCogMessageDueToTimeoutIfTimedOut(requests.CogRedisMessage{
			Input:  cogReqBody.Input,
			Error:  "TIMEOUT",
			Status: requests.CogFailed,
		})
	}()

	render.Status(r, http.StatusOK)
	render.JSON(w, r, &responses.TaskQueuedResponse{
		ID: requestId,
	})
}
