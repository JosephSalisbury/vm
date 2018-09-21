package google

func (p *googleProvider) Delete(id string) error {
	if _, err := p.gcloud(gCloudCommand{
		description: "delete instance",
		args: []string{
			"compute",
			"instances",
			"delete",
			id,
			"--zone", zone,
		},
	}); err != nil {
		return err
	}

	return nil
}
