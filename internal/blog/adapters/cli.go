package adapters

import (
	"log"

	"github.com/edlingao/internal/auth/ports"
	"github.com/edlingao/internal/blog/core"
)

type CLIService struct {
	blogRepo    *BlogRepo
	userSerivce ports.UserService
}

func NewCLIService(
	blogRepo *BlogRepo,
	userService ports.UserService,
) *CLIService {
	return &CLIService{
		blogRepo:    blogRepo,
		userSerivce: userService,
	}
}

func (c *CLIService) SaveEntry(title string) error {
	err := c.UpdateEntry(title)
	if err == nil {
		return nil
	}

	blogObj := core.NewBlog(title)
	err = blogObj.ProcessFileAndSave()
	if err != nil {
		log.Println("Error saving blog: ", err)
		return nil
	}

	_, err = c.blogRepo.Save(blogObj)
	if err != nil {
		log.Println("Error saving blog: ", err)
		return err
	}

	err = c.blogRepo.AddTagsToBlog(blogObj.ID, blogObj.GetTags())
	if err != nil {
		log.Println("Error adding tags to blog: ", err)
		return err
	}

	return err
}

func (c *CLIService) UpdateEntry(title string) error {
	blog, err := c.blogRepo.GetByTitle(title)
	if err != nil {
		return err
	}

	err = blog.ProcessFileAndSave()
	if err != nil {
		log.Println("Error saving blog during update: ", err)
		return err
	}

	_, err = c.blogRepo.Update(blog)
	if err != nil {
		log.Println("Error updating blog: ", err)
		return err
	}

	err = c.blogRepo.AddTagsToBlog(blog.ID, blog.GetTags())
	if err != nil {
		log.Println("Error updating blog tags: ", err)
		return err
	}

	return nil
}

func (c *CLIService) AddUser(username, password, role string) error {
	return c.userSerivce.Register(username, password, role)
}
