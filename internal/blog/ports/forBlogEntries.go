package ports

type BlogCLISerivce interface {
	SaveEntry(title string) error
}
