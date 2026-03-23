package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Paste holds the schema definition for the Paste entity.
type Paste struct {
	ent.Schema
}

// Fields of the Paste.
func (Paste) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("id").
			Positive(),

		field.Int64("owner_id").
			Positive().
			Comment("关联的用户 ID"),

		field.String("title").
			NotEmpty().
			MaxLen(120).
			Comment("代码片段标题"),

		field.String("short_link").
			Optional().
			Unique().
			MaxLen(10).
			Comment("短链接标识"),

		field.Text("content").
			NotEmpty().
			Comment("代码或文本内容"),

		field.String("language").
			Default("text").
			MaxLen(20).
			Comment("编程语言"),

		field.Enum("visibility").
			Values("public", "private").
			Default("private").
			Comment("可见性"),

		field.Time("created_at").
			Default(time.Now).
			Immutable().
			Comment("创建时间"),

		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Comment("最后更新时间"),
	}
}

// Edges of the Paste.
func (Paste) Edges() []ent.Edge {
	return nil
}

// Indexes of the Paste.
func (Paste) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("owner_id"),
	}
}
