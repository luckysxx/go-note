package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
)

// Snippet 代码片段/笔记/文件的统一实体，替代原有的旧内容模型。
// 通过 type 字段区分内容类型：code（代码片段）、note（笔记）、file（文件上传）。
type Snippet struct {
	ent.Schema
}

// Fields of the Snippet.
func (Snippet) Fields() []ent.Field {
	return []ent.Field{
		// ── 主键：由 id-generator 雪花算法生成 ──
		field.Int64("id").
			Comment("雪花算法生成的全局唯一 ID"),

		field.Int64("owner_id").
			Positive().
			Comment("关联的用户 ID"),

		// ── 内容类型区分 ──
		field.Enum("type").
			Values("code", "note", "file").
			Default("code").
			Comment("内容类型：code=代码片段 note=笔记 file=文件上传"),

		field.String("title").
			NotEmpty().
			MaxLen(200).
			Comment("标题"),

		// ── 文本内容（code/note 类型使用）──
		field.Text("content").
			Optional().
			Comment("代码或笔记的文本内容"),

		// ── 文件信息（file 类型使用）──
		field.String("file_url").
			Optional().
			MaxLen(512).
			Comment("文件的对象存储路径"),

		field.Int64("file_size").
			Optional().
			Default(0).
			Comment("文件大小（字节）"),

		field.String("mime_type").
			Optional().
			MaxLen(100).
			Comment("文件 MIME 类型"),

		// ── 编程语言 ──
		field.String("language").
			Default("text").
			MaxLen(30).
			Comment("编程语言或文件类型标识"),

		// ── 可见性 ──
		field.Enum("visibility").
			Values("public", "private").
			Default("private").
			Comment("可见性：public=公开 private=私有"),

		// ── 分组关联 ──
		field.Int64("group_id").
			Optional().
			Nillable().
			Comment("所属分组 ID"),

		// ── 时间戳 ──
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

// Edges of the Snippet.
func (Snippet) Edges() []ent.Edge {
	return []ent.Edge{
		// Snippet → Group（多对一）
		edge.From("group", Group.Type).
			Ref("snippets").
			Field("group_id").
			Unique(),

		// Snippet ↔ Tag（多对多）
		edge.To("tags", Tag.Type),
	}
}

// Indexes of the Snippet.
func (Snippet) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("owner_id"),
		index.Fields("owner_id", "type"),
		index.Fields("owner_id", "group_id"),
		index.Fields("visibility"),
	}
}
