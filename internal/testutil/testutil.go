package testutil

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/PacketFabric/terraform-provider-packetfabric/internal/packetfabric"
	"github.com/google/uuid"
)

var birdNames = []string{
	"albatross",
	"blackbird",
	"canary",
	"dove",
	"eagle",
	"falcon",
	"goldfinch",
	"hawk",
	"ibis",
	"jay",
	"kite",
	"lark",
	"magpie",
	"nightingale",
	"owl",
	"parrot",
	"quail",
	"raven",
	"sparrow",
	"toucan",
	"vulture",
	"woodpecker",
	"xantus",
	"yellowhammer",
	"zebrafinch",
	"ape",
	"bear",
	"cheetah",
	"dog",
	"elephant",
	"fox",
	"giraffe",
	"hippo",
	"iguana",
	"jaguar",
	"koala",
	"lion",
	"monkey",
	"narwhal",
	"otter",
	"panda",
	"quokka",
	"rabbit",
	"sheep",
	"tiger",
	"unicorn",
	"viper",
	"whale",
	"yak",
	"zebra",
}

func GenerateUniqueName() string {
	rand.Seed(time.Now().UnixNano())
	birdName := birdNames[rand.Intn(len(birdNames))]
	return fmt.Sprintf("terraform_testacc_%s", birdName)
}

func GenerateUniqueResourceName(resource string) (resourceName, hclName string) {
	uuid := uuid.NewString()
	shortUuid := uuid[0:8]
	randomNumber := rand.Intn(9000) + 1000
	hclName = fmt.Sprintf("terraform_testacc_%s_%d", shortUuid, randomNumber)
	resourceName = fmt.Sprintf("%s.%s", resource, hclName)
	return
}

func _createPFClient() (*packetfabric.PFClient, error) {
	host := os.Getenv("PF_HOST")
	token := os.Getenv("PF_TOKEN")
	c, err := packetfabric.NewPFClient(&host, &token)
	if err != nil {
		return nil, fmt.Errorf("error creating PFClient: %w", err)
	}
	return c, nil
}

func PreCheck(t *testing.T, additionalEnvVars []string) {
	requiredEnvVars := []string{
		"PF_HOST",
		"PF_TOKEN",
		"PF_ACCOUNT_ID",
	}
	if additionalEnvVars != nil {
		requiredEnvVars = append(requiredEnvVars, additionalEnvVars...)
	}
	missing := false
	for _, variable := range requiredEnvVars {
		if _, ok := os.LookupEnv(variable); !ok {
			missing = true
			t.Errorf("`%s` must be set for this acceptance test!", variable)
		}
	}
	if missing {
		t.Fatalf("Some environment variables missing.")
	}
}

func _contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

// Used for Port, Point-to-Point, Cloud Router IPsec
func GetPopAndZoneWithAvailablePort(desiredSpeed string, SkipDesiredMarket *string, DesiredMedia *string, IsIpsecCapable bool) (pop, zone, media, market string, availabilityErr error) {

	c, err := _createPFClient()
	if err != nil {
		log.Println("Error creating PF client: ", err)
		return "", "", "", "", err
	}

	var locations []packetfabric.Location
	if IsIpsecCapable {
		log.Println("IPsec capable set to true")
		locations, err = c.ListLocationsIpsecCapable()
		if err != nil {
			log.Println("Error getting locations list: ", err)
			return "", "", "", "", fmt.Errorf("error getting locations list: %w", err)
		}
	} else {
		locations, err = c.ListLocations()
		if err != nil {
			log.Println("Error getting locations list: ", err)
			return "", "", "", "", fmt.Errorf("error getting locations list: %w", err)
		}
	}

	// We need to shuffle the list of locations. Otherwise, we may try to run
	// all tests on the same port.
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(locations), func(i, j int) { locations[i], locations[j] = locations[j], locations[i] })

	testingInLab := strings.Contains(os.Getenv("PF_HOST"), "api.dev")

	for _, l := range locations {
		// log.Printf("Checking PoP: %s\n", l.Pop)
		// Skip Colt locations
		if l.Vendor == "Colt" {
			continue
		}

		// Do not select a port in the same market as the one set in SkipDesiredMarket
		if SkipDesiredMarket != nil && l.Market == *SkipDesiredMarket {
			continue
		}
		portAvailability, err := c.GetLocationPortAvailability(l.Pop)
		if err != nil {
			log.Println("Error getting location port availability for ", l.Pop, ": ", err)
			return "", "", "", "", fmt.Errorf("error getting location port availability: %w", err)
		}

		for _, p := range portAvailability {
			// if DesiredMedia specified, only select ports with that media
			if DesiredMedia != nil {
				if p.Speed == desiredSpeed && p.Media == *DesiredMedia && p.Count > 0 && (!testingInLab || _contains(labPopsPort, l.Pop)) {
					pop = l.Pop
					zone = p.Zone
					media = p.Media
					market = l.Market
					log.Println("Specified Media to match: ", *DesiredMedia)
					log.Println("Found available port at ", pop, zone, media, market)
					return
				}
			} else {
				if p.Speed == desiredSpeed && p.Count > 0 && (!testingInLab || _contains(labPopsPort, l.Pop)) {
					pop = l.Pop
					zone = p.Zone
					media = p.Media
					market = l.Market
					log.Println("Found available port at ", pop, zone, media, market)
					if SkipDesiredMarket == nil {
						log.Println("Not specified Market to avoid.")
					} else {
						log.Println("Specified Market to avoid: ", *SkipDesiredMarket)
					}
					return
				}
			}
		}
	}
	log.Println("No pops with available ports found.")
	return "", "", "", "", errors.New("no pops with available ports")
}

// Used for Hosted Cloud and Cloud Router Connection
func (details PortDetails) FindAvailableCloudPopZone() (pop, zone, region string) {
	popsWithZones, _ := details.FetchCloudPopsAndZones()
	popsToSkip := make([]string, 0)

	log.Println("Starting to search for available Cloud PoP and zone...")
	log.Printf("Available PoPs with Zones: %v\n", popsWithZones)

	testingInLab := strings.Contains(os.Getenv("PF_HOST"), "api.dev")

	// First loop, prioritizing "us-" regions
	for popAvailable, zones := range popsWithZones {
		region = zones[len(zones)-1] // always take the last zone as region
		if len(zones) > 1 && (!testingInLab || _contains(labPopsHostedCloud, popAvailable)) && strings.HasPrefix(region, "us-") {
			pop = popAvailable
			zone = zones[0] // always take the first zone available
			log.Printf("Found available Hosted Cloud PoP: %s, Zone: %s, Region: %s\n", pop, zone, region)
			return
		} else {
			popsToSkip = append(popsToSkip, popAvailable)
		}
	}

	// Second loop, if no "us-" region PoP is found
	for popAvailable, zones := range popsWithZones {
		if len(popsToSkip) == len(popsWithZones) {
			log.Fatal(errors.New("there's no port available on any pop"))
		}
		if _contains(popsToSkip, popAvailable) {
			log.Printf("PoP %s is in popsToSkip, skipping...\n", popAvailable)
			continue
		} else {
			if len(zones) > 1 && (!testingInLab || _contains(labPopsHostedCloud, popAvailable)) {
				pop = popAvailable
				zone = zones[0] // always take the first zone available
				region = zones[len(zones)-1]
				log.Printf("Found available Hosted Cloud PoP: %s, Zone: %s, Region: %s\n", pop, zone, region)
				return
			} else {
				popsToSkip = append(popsToSkip, popAvailable)
			}
		}
	}

	log.Println("No available Hosted Cloud PoP, zone, and region found.")
	return
}

// Used for Dedicated Cloud
func (details PortDetails) FindAvailableCloudPopZoneDedicated() (pop, zone, region string, availabilityErr error) {
	popsWithZones, _ := details.FetchCloudPopsAndZones()
	popsToSkip := make([]string, 0)

	c, err := _createPFClient()
	if err != nil {
		log.Println("Error creating PF client: ", err)
		return "", "", "", err
	}

	log.Println("Starting to search for available Cloud PoP and zone...")
	log.Printf("Available PoPs with Zones: %v\n", popsWithZones)

	testingInLab := strings.Contains(os.Getenv("PF_HOST"), "api.dev")

	for popAvailable, zones := range popsWithZones {
		// log.Printf("Checking PoP: %s\n", popAvailable)
		if len(popsToSkip) == len(popsWithZones) {
			log.Fatal(errors.New("there's no port available on any pop"))
		}
		if _contains(popsToSkip, popAvailable) {
			log.Printf("PoP %s is in popsToSkip, skipping...\n", popAvailable)
			continue
		} else {
			if len(zones) > 1 && (!testingInLab || _contains(labPopsDedicatedCloud, popAvailable)) {
				portAvailability, err := c.GetLocationPortAvailability(popAvailable)
				if err != nil {
					log.Println("Error getting location port availability for ", popAvailable, ": ", err)
					return "", "", "", fmt.Errorf("error getting location port availability: %w", err)
				}
				for _, p := range portAvailability {
					if p.Speed == details.DesiredSpeed && p.Count > 0 {
						pop = popAvailable
						zone = p.Zone // get the zone available via /v2/locations/%s/port-availability
						region = zones[len(zones)-1]
						log.Printf("Found available Dedicated Cloud PoP: %s, Zone: %s, Region: %s\n", pop, zone, region)
						return
					}
				}
				return
			} else {
				popsToSkip = append(popsToSkip, popAvailable)
			}
		}
	}
	log.Println("No available Dedicated Cloud PoP, zone found.")
	return
}

func (details PortDetails) FetchCloudPopsAndZones() (popsWithZones map[string][]string, err error) {
	if details.DesiredProvider == "" {
		err = errors.New("please provide a valid cloud provider to fetch pop")
	}
	if details.PFClient == nil {
		err = errors.New("please create PFClient to fetch cloud pop")
		return
	}
	popsWithZones = make(map[string][]string)
	if cloudLocations, locErr := details.PFClient.GetCloudLocations(
		details.DesiredProvider,
		details.DesiredConnectionType,
		details.IsNatCapable,
		details.HasCloudRouter,
		details.AnyType,
		details.DesiredPop,
		details.DesiredCity,
		details.DesiredState,
		details.DesiredMarket,
		details.DesiredRegion); locErr != nil {
		err = locErr
		return
	} else {
		for _, loc := range cloudLocations {
			popsWithZones[loc.Pop] = append(loc.Zones, loc.CloudConnectionDetails.Region)

		}
	}
	return
}

func CreateBasePortDetails() PortDetails {
	c, err := _createPFClient()
	if err != nil {
		log.Panic(err)
	}
	return PortDetails{
		PFClient:          c,
		DesiredSpeed:      portSpeed,
		SkipDesiredMarket: nil,
	}
}

func setAzureLocations(host string) (string, string, string) {
	var AzureLocation string
	var AzurePeeringLocation string
	var AzureServiceProviderName string

	testingInLab := strings.Contains(host, "api.dev")

	if testingInLab {
		AzureLocation = AzureLocationDev
		AzurePeeringLocation = AzurePeeringLocationDev
		AzureServiceProviderName = AzureServiceProviderNameDev
	} else {
		AzureLocation = AzureLocationProd
		AzurePeeringLocation = AzurePeeringLocationProd
		AzureServiceProviderName = AzureServiceProviderNameProd
	}

	return AzureLocation, AzurePeeringLocation, AzureServiceProviderName
}
