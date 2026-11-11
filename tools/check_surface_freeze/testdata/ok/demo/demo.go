package demo

import "context"

// Kind 是夹具里的具名类型。
type Kind string

const (
	KindA Kind = "a"
	KindB Kind = "b"
)

// DefaultName 是无显式类型的常量：指纹冻结它的字面量。
const DefaultName = "demo"

// Widget 是导出结构体：字段名与类型都进指纹。
type Widget struct {
	Name string
	Size int
}

// Label 是值接收者方法。
func (w Widget) Label() string { return w.Name }

// Resize 是指针接收者方法：接口满足关系不同，指纹要能区分。
func (w *Widget) Resize(n int) error { return nil }

// Store 是导出接口：方法集进指纹。
type Store interface {
	Get(ctx context.Context, id string) (Widget, error)
}

// New 是带可变参数与函数选项的构造函数。
func New(name string, opts ...func(*Widget)) *Widget {
	w := &Widget{Name: name}
	for _, opt := range opts {
		opt(w)
	}
	return w
}
