package domain

// Phase represents a lifecycle phase in IT project documentation.
type Phase string

const (
	PhaseInitiation    Phase = "initiation"
	PhasePlanning      Phase = "planning"
	PhaseDesign        Phase = "design"
	PhaseDevelopment   Phase = "development"
	PhaseTesting       Phase = "testing"
	PhaseDeployment    Phase = "deployment"
	PhaseMaintenance   Phase = "maintenance"
	PhaseDecommission  Phase = "decommission"
	PhaseUnspecified   Phase = "unspecified"
)

// AllPhases returns phases in lifecycle order.
func AllPhases() []Phase {
	return []Phase{
		PhaseInitiation,
		PhasePlanning,
		PhaseDesign,
		PhaseDevelopment,
		PhaseTesting,
		PhaseDeployment,
		PhaseMaintenance,
		PhaseDecommission,
	}
}
