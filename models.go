package blog

import (
  "fmt"
  "net/url"
  "strings"
  "time"

  "appengine"
  "appengine/datastore"
)

type Post struct {
  Slug *datastore.Key
  Title string
  Content string
  PublishedDate time.Time
  EditDate time.Time
  IsDraft bool
}

func (p *Post) Url() *url.URL {
  url, err := r.Get("post").URL("date", p.PublishedDate.Format("2006/01/02"), "slug", p.Slug.StringID())
  if err != nil {
    panic(err)
  }
  return url
}

func (p *Post) EditUrl() *url.URL {
  url, err := r.Get("EditPost").URL("slug", p.Slug.StringID())
  if err != nil {
    panic(err)
  }
  return url
}

func hyphenizeTitle (title string) string {
  return strings.Replace(title, " ", "-", -1)
}

func postKey(c appengine.Context, slug string) *datastore.Key {
  return datastore.NewKey(c, "Posts", slug, 0, nil)
}

func createSlug(ctx appengine.Context, p *Post) *datastore.Key {
  if p.Slug != nil {
    return p.Slug
  }
  hyphenatedTitle := hyphenizeTitle(p.Title)
  return postKey(ctx, hyphenatedTitle)
}

func newPost(c appengine.Context) (*Post, error) {
  post := &Post{
    Title: fmt.Sprintf("Post %s", time.Now().Format("20060102030405")),
    Content: "Fill me in",
    PublishedDate: time.Now(),
    EditDate: time.Now(),
    IsDraft: true,
  }
  slug := createSlug(c, post)
  post.Slug = slug

  _, err := datastore.Put(c, slug, post)
  return post, err
}

func getPostFromSlugString(c appengine.Context, slugString string) (*Post, error) {
  key := postKey(c, slugString)
  post := &Post{Slug: key}
  err := datastore.Get(c, key, post)
  return post, err
}

func getAllPosts(c appengine.Context) ([]Post, error) {
  count, countErr := datastore.NewQuery("Posts").Count(c)
  if countErr != nil {
    posts := make([]Post, 0)
    return posts, countErr
  }

  posts := make([]Post, 0, count)
  _, getErr := datastore.NewQuery("Posts").Order("-PublishedDate").GetAll(c, &posts)
  return posts, getErr
}

func updatePost(c appengine.Context, slugString string, updatePost *Post) (*Post, error) {
  post, err := getPostFromSlugString(c, slugString)
  if err != nil {
    return post, err
  }

  // Compare the stored post to the updated value
  if post.Title != updatePost.Title {
    key := postKey(c, slugString)
    datastore.Delete(c, key)
    post.Slug = nil
    post.Title = updatePost.Title
    post.Slug = createSlug(c, post)
  }

  post.Content = updatePost.Content
  post.IsDraft = updatePost.IsDraft
  _, puterr := datastore.Put(c, post.Slug, post)
  return post, puterr
}