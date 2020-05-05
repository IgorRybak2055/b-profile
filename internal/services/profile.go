package services

import (
	"context"
	"crypto/rand"
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/IgorRybak2055/bamboo/internal/models"
	"github.com/IgorRybak2055/bamboo/internal/repository"
)

var ErrGenerateMatchingID = errors.New("cannot generate matchingID")

type profileService struct {
	profileRepo repository.Profile
	log         *logrus.Logger
}

// NewAccountService will create new accountService object representation of Account interface
func NewAccountService(pr repository.Profile, logger *logrus.Logger) Profile {
	return &profileService{
		profileRepo: pr,
		log:         logger,
	}
}

func generateMatchingID() (string, error) {
	b := make([]byte, 12)

	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	matchingID := strings.ToUpper(fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]))

	return matchingID, nil
}

func (p profileService) OrderProfile(ctx context.Context, iccids []string) ([][]string, error) {
	var (
		csvHeader = []string{"iccid", "matchingId", "error"}
		profiles  = make([][]string, 0, len(iccids))
		ch        = make(chan []string, len(iccids))
	)

	profiles = append(profiles, csvHeader)

	for _, iccid := range iccids {
		go p.orderWorker(ctx, iccid, ch)
	}

	for len(profiles) < len(iccids) {
		profiles = append(profiles, <-ch)
	}

	return profiles, nil
}

func (p profileService) orderWorker(ctx context.Context, iccid string, ch chan []string) {
	var (
		profile = models.Profile{ICCID: iccid}
		err     error
	)

	if profile.MatchingID, err = generateMatchingID(); err != nil {
		p.log.Warnf("failed to generate matchingId for %s iccid: %s", iccid, err)

		profile.MatchingID = ""
		data := profile.CVSRespond()
		data = append(data, ErrGenerateMatchingID.Error())

		ch <- data

		return
	}

	if err = p.profileRepo.OrderProfile(ctx, profile); err != nil {
		p.log.Warnf("failed to save matchingId %s iccid: %s", iccid, err)

		profile.MatchingID = ""

		data := profile.CVSRespond()
		data = append(data, err.Error())

		ch <- data

		return
	}

	data := profile.CVSRespond()
	data = append(data, "")

	ch <- data
}
