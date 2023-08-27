package try

type Try interface {
	Run() error
}

func TryN(t *Try, n int) error {
	var err error
	for i := 0; i < n; i++ {
		err = (*t).Run()
		if err != nil {
			continue
		}
		return nil
	}
	return err
}
