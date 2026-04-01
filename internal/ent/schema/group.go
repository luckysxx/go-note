package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Group 分组实体，用于组织和分类 Snippet。
type Group struct {
	ent.Schema
}

// Fields of the Group.
func (Group) Fields() []ent.Field {
	return []ent.Field{
		// 主键：由 id-generator 雪花算法生成
		field.Int64("id").
			Comment("雪花算法生成的全局唯一 ID"),

		field.Int64("owner_id").
			Positive().
			Comment("分组拥有者的用户 ID"),

		field.String("name").
			NotEmpty().
			MaxLen(60).
			Comment("分组名称"),

		field.Int("sort_order").
			Default(0).
			Comment("排序权重，值越小越靠前"),

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

// Edges of the Group.
func (Group) Edges() []ent.Edge {
	return []ent.Edge{
		// Group → Snippet（一对多）
		edge.To("snippets", Snippet.Type),
	}
}

// Indexes of the Group.
func (Group) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("owner_id"),
		index.Fields("owner_id", "name").Unique(),
	}
}
