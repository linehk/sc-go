// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/stablecog/sc-go/database/ent/generation"
	"github.com/stablecog/sc-go/database/ent/generationoutput"
	"github.com/stablecog/sc-go/database/ent/predicate"
)

// GenerationOutputUpdate is the builder for updating GenerationOutput entities.
type GenerationOutputUpdate struct {
	config
	hooks     []Hook
	mutation  *GenerationOutputMutation
	modifiers []func(*sql.UpdateBuilder)
}

// Where appends a list predicates to the GenerationOutputUpdate builder.
func (gou *GenerationOutputUpdate) Where(ps ...predicate.GenerationOutput) *GenerationOutputUpdate {
	gou.mutation.Where(ps...)
	return gou
}

// SetImagePath sets the "image_path" field.
func (gou *GenerationOutputUpdate) SetImagePath(s string) *GenerationOutputUpdate {
	gou.mutation.SetImagePath(s)
	return gou
}

// SetUpscaledImagePath sets the "upscaled_image_path" field.
func (gou *GenerationOutputUpdate) SetUpscaledImagePath(s string) *GenerationOutputUpdate {
	gou.mutation.SetUpscaledImagePath(s)
	return gou
}

// SetNillableUpscaledImagePath sets the "upscaled_image_path" field if the given value is not nil.
func (gou *GenerationOutputUpdate) SetNillableUpscaledImagePath(s *string) *GenerationOutputUpdate {
	if s != nil {
		gou.SetUpscaledImagePath(*s)
	}
	return gou
}

// ClearUpscaledImagePath clears the value of the "upscaled_image_path" field.
func (gou *GenerationOutputUpdate) ClearUpscaledImagePath() *GenerationOutputUpdate {
	gou.mutation.ClearUpscaledImagePath()
	return gou
}

// SetGalleryStatus sets the "gallery_status" field.
func (gou *GenerationOutputUpdate) SetGalleryStatus(gs generationoutput.GalleryStatus) *GenerationOutputUpdate {
	gou.mutation.SetGalleryStatus(gs)
	return gou
}

// SetNillableGalleryStatus sets the "gallery_status" field if the given value is not nil.
func (gou *GenerationOutputUpdate) SetNillableGalleryStatus(gs *generationoutput.GalleryStatus) *GenerationOutputUpdate {
	if gs != nil {
		gou.SetGalleryStatus(*gs)
	}
	return gou
}

// SetGenerationID sets the "generation_id" field.
func (gou *GenerationOutputUpdate) SetGenerationID(u uuid.UUID) *GenerationOutputUpdate {
	gou.mutation.SetGenerationID(u)
	return gou
}

// SetDeletedAt sets the "deleted_at" field.
func (gou *GenerationOutputUpdate) SetDeletedAt(t time.Time) *GenerationOutputUpdate {
	gou.mutation.SetDeletedAt(t)
	return gou
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (gou *GenerationOutputUpdate) SetNillableDeletedAt(t *time.Time) *GenerationOutputUpdate {
	if t != nil {
		gou.SetDeletedAt(*t)
	}
	return gou
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (gou *GenerationOutputUpdate) ClearDeletedAt() *GenerationOutputUpdate {
	gou.mutation.ClearDeletedAt()
	return gou
}

// SetUpdatedAt sets the "updated_at" field.
func (gou *GenerationOutputUpdate) SetUpdatedAt(t time.Time) *GenerationOutputUpdate {
	gou.mutation.SetUpdatedAt(t)
	return gou
}

// SetGenerationsID sets the "generations" edge to the Generation entity by ID.
func (gou *GenerationOutputUpdate) SetGenerationsID(id uuid.UUID) *GenerationOutputUpdate {
	gou.mutation.SetGenerationsID(id)
	return gou
}

// SetGenerations sets the "generations" edge to the Generation entity.
func (gou *GenerationOutputUpdate) SetGenerations(g *Generation) *GenerationOutputUpdate {
	return gou.SetGenerationsID(g.ID)
}

// Mutation returns the GenerationOutputMutation object of the builder.
func (gou *GenerationOutputUpdate) Mutation() *GenerationOutputMutation {
	return gou.mutation
}

// ClearGenerations clears the "generations" edge to the Generation entity.
func (gou *GenerationOutputUpdate) ClearGenerations() *GenerationOutputUpdate {
	gou.mutation.ClearGenerations()
	return gou
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (gou *GenerationOutputUpdate) Save(ctx context.Context) (int, error) {
	gou.defaults()
	return withHooks[int, GenerationOutputMutation](ctx, gou.sqlSave, gou.mutation, gou.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (gou *GenerationOutputUpdate) SaveX(ctx context.Context) int {
	affected, err := gou.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (gou *GenerationOutputUpdate) Exec(ctx context.Context) error {
	_, err := gou.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (gou *GenerationOutputUpdate) ExecX(ctx context.Context) {
	if err := gou.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (gou *GenerationOutputUpdate) defaults() {
	if _, ok := gou.mutation.UpdatedAt(); !ok {
		v := generationoutput.UpdateDefaultUpdatedAt()
		gou.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (gou *GenerationOutputUpdate) check() error {
	if v, ok := gou.mutation.GalleryStatus(); ok {
		if err := generationoutput.GalleryStatusValidator(v); err != nil {
			return &ValidationError{Name: "gallery_status", err: fmt.Errorf(`ent: validator failed for field "GenerationOutput.gallery_status": %w`, err)}
		}
	}
	if _, ok := gou.mutation.GenerationsID(); gou.mutation.GenerationsCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "GenerationOutput.generations"`)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (gou *GenerationOutputUpdate) Modify(modifiers ...func(u *sql.UpdateBuilder)) *GenerationOutputUpdate {
	gou.modifiers = append(gou.modifiers, modifiers...)
	return gou
}

func (gou *GenerationOutputUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := gou.check(); err != nil {
		return n, err
	}
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   generationoutput.Table,
			Columns: generationoutput.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUUID,
				Column: generationoutput.FieldID,
			},
		},
	}
	if ps := gou.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := gou.mutation.ImagePath(); ok {
		_spec.SetField(generationoutput.FieldImagePath, field.TypeString, value)
	}
	if value, ok := gou.mutation.UpscaledImagePath(); ok {
		_spec.SetField(generationoutput.FieldUpscaledImagePath, field.TypeString, value)
	}
	if gou.mutation.UpscaledImagePathCleared() {
		_spec.ClearField(generationoutput.FieldUpscaledImagePath, field.TypeString)
	}
	if value, ok := gou.mutation.GalleryStatus(); ok {
		_spec.SetField(generationoutput.FieldGalleryStatus, field.TypeEnum, value)
	}
	if value, ok := gou.mutation.DeletedAt(); ok {
		_spec.SetField(generationoutput.FieldDeletedAt, field.TypeTime, value)
	}
	if gou.mutation.DeletedAtCleared() {
		_spec.ClearField(generationoutput.FieldDeletedAt, field.TypeTime)
	}
	if value, ok := gou.mutation.UpdatedAt(); ok {
		_spec.SetField(generationoutput.FieldUpdatedAt, field.TypeTime, value)
	}
	if gou.mutation.GenerationsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   generationoutput.GenerationsTable,
			Columns: []string{generationoutput.GenerationsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUUID,
					Column: generation.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := gou.mutation.GenerationsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   generationoutput.GenerationsTable,
			Columns: []string{generationoutput.GenerationsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUUID,
					Column: generation.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_spec.AddModifiers(gou.modifiers...)
	if n, err = sqlgraph.UpdateNodes(ctx, gou.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{generationoutput.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	gou.mutation.done = true
	return n, nil
}

// GenerationOutputUpdateOne is the builder for updating a single GenerationOutput entity.
type GenerationOutputUpdateOne struct {
	config
	fields    []string
	hooks     []Hook
	mutation  *GenerationOutputMutation
	modifiers []func(*sql.UpdateBuilder)
}

// SetImagePath sets the "image_path" field.
func (gouo *GenerationOutputUpdateOne) SetImagePath(s string) *GenerationOutputUpdateOne {
	gouo.mutation.SetImagePath(s)
	return gouo
}

// SetUpscaledImagePath sets the "upscaled_image_path" field.
func (gouo *GenerationOutputUpdateOne) SetUpscaledImagePath(s string) *GenerationOutputUpdateOne {
	gouo.mutation.SetUpscaledImagePath(s)
	return gouo
}

// SetNillableUpscaledImagePath sets the "upscaled_image_path" field if the given value is not nil.
func (gouo *GenerationOutputUpdateOne) SetNillableUpscaledImagePath(s *string) *GenerationOutputUpdateOne {
	if s != nil {
		gouo.SetUpscaledImagePath(*s)
	}
	return gouo
}

// ClearUpscaledImagePath clears the value of the "upscaled_image_path" field.
func (gouo *GenerationOutputUpdateOne) ClearUpscaledImagePath() *GenerationOutputUpdateOne {
	gouo.mutation.ClearUpscaledImagePath()
	return gouo
}

// SetGalleryStatus sets the "gallery_status" field.
func (gouo *GenerationOutputUpdateOne) SetGalleryStatus(gs generationoutput.GalleryStatus) *GenerationOutputUpdateOne {
	gouo.mutation.SetGalleryStatus(gs)
	return gouo
}

// SetNillableGalleryStatus sets the "gallery_status" field if the given value is not nil.
func (gouo *GenerationOutputUpdateOne) SetNillableGalleryStatus(gs *generationoutput.GalleryStatus) *GenerationOutputUpdateOne {
	if gs != nil {
		gouo.SetGalleryStatus(*gs)
	}
	return gouo
}

// SetGenerationID sets the "generation_id" field.
func (gouo *GenerationOutputUpdateOne) SetGenerationID(u uuid.UUID) *GenerationOutputUpdateOne {
	gouo.mutation.SetGenerationID(u)
	return gouo
}

// SetDeletedAt sets the "deleted_at" field.
func (gouo *GenerationOutputUpdateOne) SetDeletedAt(t time.Time) *GenerationOutputUpdateOne {
	gouo.mutation.SetDeletedAt(t)
	return gouo
}

// SetNillableDeletedAt sets the "deleted_at" field if the given value is not nil.
func (gouo *GenerationOutputUpdateOne) SetNillableDeletedAt(t *time.Time) *GenerationOutputUpdateOne {
	if t != nil {
		gouo.SetDeletedAt(*t)
	}
	return gouo
}

// ClearDeletedAt clears the value of the "deleted_at" field.
func (gouo *GenerationOutputUpdateOne) ClearDeletedAt() *GenerationOutputUpdateOne {
	gouo.mutation.ClearDeletedAt()
	return gouo
}

// SetUpdatedAt sets the "updated_at" field.
func (gouo *GenerationOutputUpdateOne) SetUpdatedAt(t time.Time) *GenerationOutputUpdateOne {
	gouo.mutation.SetUpdatedAt(t)
	return gouo
}

// SetGenerationsID sets the "generations" edge to the Generation entity by ID.
func (gouo *GenerationOutputUpdateOne) SetGenerationsID(id uuid.UUID) *GenerationOutputUpdateOne {
	gouo.mutation.SetGenerationsID(id)
	return gouo
}

// SetGenerations sets the "generations" edge to the Generation entity.
func (gouo *GenerationOutputUpdateOne) SetGenerations(g *Generation) *GenerationOutputUpdateOne {
	return gouo.SetGenerationsID(g.ID)
}

// Mutation returns the GenerationOutputMutation object of the builder.
func (gouo *GenerationOutputUpdateOne) Mutation() *GenerationOutputMutation {
	return gouo.mutation
}

// ClearGenerations clears the "generations" edge to the Generation entity.
func (gouo *GenerationOutputUpdateOne) ClearGenerations() *GenerationOutputUpdateOne {
	gouo.mutation.ClearGenerations()
	return gouo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (gouo *GenerationOutputUpdateOne) Select(field string, fields ...string) *GenerationOutputUpdateOne {
	gouo.fields = append([]string{field}, fields...)
	return gouo
}

// Save executes the query and returns the updated GenerationOutput entity.
func (gouo *GenerationOutputUpdateOne) Save(ctx context.Context) (*GenerationOutput, error) {
	gouo.defaults()
	return withHooks[*GenerationOutput, GenerationOutputMutation](ctx, gouo.sqlSave, gouo.mutation, gouo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (gouo *GenerationOutputUpdateOne) SaveX(ctx context.Context) *GenerationOutput {
	node, err := gouo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (gouo *GenerationOutputUpdateOne) Exec(ctx context.Context) error {
	_, err := gouo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (gouo *GenerationOutputUpdateOne) ExecX(ctx context.Context) {
	if err := gouo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (gouo *GenerationOutputUpdateOne) defaults() {
	if _, ok := gouo.mutation.UpdatedAt(); !ok {
		v := generationoutput.UpdateDefaultUpdatedAt()
		gouo.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (gouo *GenerationOutputUpdateOne) check() error {
	if v, ok := gouo.mutation.GalleryStatus(); ok {
		if err := generationoutput.GalleryStatusValidator(v); err != nil {
			return &ValidationError{Name: "gallery_status", err: fmt.Errorf(`ent: validator failed for field "GenerationOutput.gallery_status": %w`, err)}
		}
	}
	if _, ok := gouo.mutation.GenerationsID(); gouo.mutation.GenerationsCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "GenerationOutput.generations"`)
	}
	return nil
}

// Modify adds a statement modifier for attaching custom logic to the UPDATE statement.
func (gouo *GenerationOutputUpdateOne) Modify(modifiers ...func(u *sql.UpdateBuilder)) *GenerationOutputUpdateOne {
	gouo.modifiers = append(gouo.modifiers, modifiers...)
	return gouo
}

func (gouo *GenerationOutputUpdateOne) sqlSave(ctx context.Context) (_node *GenerationOutput, err error) {
	if err := gouo.check(); err != nil {
		return _node, err
	}
	_spec := &sqlgraph.UpdateSpec{
		Node: &sqlgraph.NodeSpec{
			Table:   generationoutput.Table,
			Columns: generationoutput.Columns,
			ID: &sqlgraph.FieldSpec{
				Type:   field.TypeUUID,
				Column: generationoutput.FieldID,
			},
		},
	}
	id, ok := gouo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "GenerationOutput.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := gouo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, generationoutput.FieldID)
		for _, f := range fields {
			if !generationoutput.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != generationoutput.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := gouo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := gouo.mutation.ImagePath(); ok {
		_spec.SetField(generationoutput.FieldImagePath, field.TypeString, value)
	}
	if value, ok := gouo.mutation.UpscaledImagePath(); ok {
		_spec.SetField(generationoutput.FieldUpscaledImagePath, field.TypeString, value)
	}
	if gouo.mutation.UpscaledImagePathCleared() {
		_spec.ClearField(generationoutput.FieldUpscaledImagePath, field.TypeString)
	}
	if value, ok := gouo.mutation.GalleryStatus(); ok {
		_spec.SetField(generationoutput.FieldGalleryStatus, field.TypeEnum, value)
	}
	if value, ok := gouo.mutation.DeletedAt(); ok {
		_spec.SetField(generationoutput.FieldDeletedAt, field.TypeTime, value)
	}
	if gouo.mutation.DeletedAtCleared() {
		_spec.ClearField(generationoutput.FieldDeletedAt, field.TypeTime)
	}
	if value, ok := gouo.mutation.UpdatedAt(); ok {
		_spec.SetField(generationoutput.FieldUpdatedAt, field.TypeTime, value)
	}
	if gouo.mutation.GenerationsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   generationoutput.GenerationsTable,
			Columns: []string{generationoutput.GenerationsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUUID,
					Column: generation.FieldID,
				},
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := gouo.mutation.GenerationsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   generationoutput.GenerationsTable,
			Columns: []string{generationoutput.GenerationsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: &sqlgraph.FieldSpec{
					Type:   field.TypeUUID,
					Column: generation.FieldID,
				},
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_spec.AddModifiers(gouo.modifiers...)
	_node = &GenerationOutput{config: gouo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, gouo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{generationoutput.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	gouo.mutation.done = true
	return _node, nil
}
