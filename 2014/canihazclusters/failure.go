package main

import "github.com/armon/consul-api"

//start block OMIT
func nodeHealthy(e *consulapi.ServiceEntry) bool {
	for _, c := range e.Checks {
		if c.Status != "passing" {
			return false
		}
	}
	return len(e.Checks) > 0
}

func serviceHealthy(e *consulapi.ServiceEntry, name string) bool {
	found := false
	for _, c := range e.Checks {
		if c.ServiceName == name {
			found = true
			if c.Status != "passing" {
				return false
			}
		}
	}
	return found
}

//end block OMIT

func main() {

}
