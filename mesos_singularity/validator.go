package mesos_singularity

import "fmt"

func validateRequestType(v interface{}, k string) (ws []string, errors []error) {
	validTypes := map[string]struct{}{
		// Since Singularity expects uppercase of these values and to make this resource
		// simpler, therefore just use upppercase.
		"SCHEDULED": {},
		"RUN_ONCE":  {},
		"SERVICE":   {},
		"WORKER":    {},
		"ON_DEMAND": {},
	}

	value := v.(string)

	if _, ok := validTypes[value]; !ok {
		errors = append(errors, fmt.Errorf(
			"%q must be one of ['SCHEDULED', 'RUN_ONCE', 'SERVICE', 'WORKER', 'ON_DEMAND']", k))
	}
	return
}

func validateRequestScheduleType(v interface{}, k string) (ws []string, errors []error) {
	// Only allow cron as valid type since this is widely known and we use this.
	// However, Singularity allows cron and quartz.
	// Since Singularity expects uppercase of these values and to make this resource
	// simpler, therefore just use upppercase.
	validTypes := map[string]struct{}{
		"CRON": {},
	}

	value := v.(string)

	if _, ok := validTypes[value]; !ok {
		errors = append(errors, fmt.Errorf(
			"%q must be only ['CRON']", k))
	}
	return
}
func validateRequestState(v interface{}, k string) (ws []string, errors []error) {
	validTypes := map[string]struct{}{
		"ACTIVE": {},
		"PAUSED": {},
	}

	value := v.(string)

	if _, ok := validTypes[value]; !ok {
		errors = append(errors, fmt.Errorf(
			"%q must be only ['ACTIVE', 'PAUSED]", k))
	}
	return
}

func validateDockerNetwork(v interface{}, k string) (ws []string, errors []error) {
	validTypes := map[string]struct{}{
		// Since Singularity expects uppercase of these values and to make this resource
		// simpler, therefore just use upppercase.
		"BRIDGE": {},
		"NONE":   {},
		"HOST":   {},
	}

	value := v.(string)

	if _, ok := validTypes[value]; !ok {
		errors = append(errors, fmt.Errorf(
			"%q must be one of ['BRIDGE', 'NONE', 'HOST']", k))
	}
	return
}

func validateSingularityPortMappingType(v interface{}, k string) (ws []string, errors []error) {
	validTypes := map[string]struct{}{
		"LITERAL":    {},
		"FROM_OFFER": {},
	}

	value := v.(string)

	if _, ok := validTypes[value]; !ok {
		errors = append(errors, fmt.Errorf(
			"%q must be one of ['LITERAL', 'FROM_OFFER']", k))
	}
	return
}

func validateSingularityPortProtocol(v interface{}, k string) (ws []string, errors []error) {
	validTypes := map[string]struct{}{
		"tcp": {},
		"udp": {},
	}

	value := v.(string)

	if _, ok := validTypes[value]; !ok {
		errors = append(errors, fmt.Errorf(
			"%q must be one of ['udp', 'tcp']", k))
	}
	return
}
