// example godocの例
package example

// User はユーザー情報を保持する構造体
type User struct {
	// ID ユーザーID
	ID string
	// Name 名前
	Name string
	// 年齢
	Age int
}

// Greet 名を名乗る
func (u *User) Greet() string {
	return "Hi! I'm " + u.Name
}
