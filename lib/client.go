package lib

import (
	"github.com/gobuffalo/buffalo"
    "golang.org/x/oauth2" 
    "github.com/digitalocean/godo" 
)

type TokenSource struct {
    AccessToken string
}

func (t *TokenSource) Token() (*oauth2.Token, error) {
    token := &oauth2.Token{
        AccessToken: t.AccessToken,
    }
    return token, nil
}

func CreateDroplet(token string, name string, region string, size string, image string, ctx buffalo.Context) (interface{}, error)  {
	// type string pat
	pat := token
	tokenSource := &TokenSource{
		AccessToken: pat,
	}
	
	oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
	client := godo.NewClient(oauthClient)

	createRequest := &godo.DropletCreateRequest{
		Name:   name,
		Region: region,
		Size:   size,
		Image: godo.DropletCreateImage{
			Slug: image,
		},
		SSHKeys: []godo.DropletCreateSSHKey{
			godo.DropletCreateSSHKey{ID: 107149},
		},
		IPv6: true,
		Tags: []string{"web"},
	}
	
	droplet, _, err := client.Droplets.Create(ctx, createRequest)
	return droplet, err
}

func FindDroplet(token string, id int, ctx buffalo.Context) (*godo.Droplet, error) {
	pat := token
	tokenSource := &TokenSource{
		AccessToken: pat,
	}
	oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
	client := godo.NewClient(oauthClient)
	droplet, _, err := client.Droplets.Get(ctx, id)
	return droplet, err
}

func ListDroplets(token string, ctx buffalo.Context) ([]godo.Droplet, error) {
	pat := token
	tokenSource := &TokenSource{
		AccessToken: pat,
	}
	oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
	client := godo.NewClient(oauthClient)

	opt := &godo.ListOptions{
		Page:    1,
		PerPage: 200,
	}
	droplets, _, err := client.Droplets.List(ctx, opt)
	return droplets, err
}

func RetriveImageBySlug(token string, slug string, ctx buffalo.Context) (*godo.Image, error) {
	pat := token
	tokenSource := &TokenSource{
		AccessToken: pat,
	}
	oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
	client := godo.NewClient(oauthClient)
	image, _, err := client.Images.GetBySlug(ctx, "ubuntu-16-04-x64")
	return image, err
}