package google

func (p *googleProvider) Delete(id string) error {
	p.logger.Printf("deleting instance")

	if _, err := p.instancesService.Delete(p.project, p.zone, id).Do(); err != nil {
		return err
	}

	p.logger.Printf("instance deleted")

	return nil
}
