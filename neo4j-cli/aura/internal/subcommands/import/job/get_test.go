package job_test

import (
	"github.com/neo4j/cli/neo4j-cli/aura/internal/test/testutils"
	"net/http"
	"testing"
)

func TestGetImportJobByIdWithoutProgress(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	mockHandler := helper.NewRequestHandlerMock("/v2beta1/projects/f607bebe-0cc0-4166-b60c-b4eed69ee7ee/import/jobs/87d485b4-73fc-4a7f-bb03-720f4672947e", http.StatusOK, `
{
    "data": {
        "id": "87d485b4-73fc-4a7f-bb03-720f4672947e",
        "import_type": "cloud",
        "info": {
            "state": "Completed",
            "start_time": "2025-08-15T13:12:51Z",
            "completion_time": "2025-08-15T13:15:19Z",
            "exit_status": {
                "state": "Cancelled",
                "message": "Cancelled"
            },
            "cancellation_requested_time": "2025-08-15T13:15:19.655209Z",
            "submitted_time": "2025-08-15T13:12:50.499797Z",
            "last_update_time": "2025-08-18T12:25:53.469873Z",
            "percentage_complete": 95.23
        },
        "data_source": {
            "id": "177216b3-49b9-4e24-a414-a9ecbb448c53",
            "type": "postgresql",
            "name": "AWS_POSTGRES_FLIGHTS"
        },
        "aura_target": {
            "db_id": "07e49cf5",
            "project_id": "f607bebe-0cc0-4166-b60c-b4eed69ee7ee"
        },
        "user_id": "f607bebe-0cc0-4166-b60c-b4eed69ee7ee"
    }
}
	`)

	helper.SetConfigValue("aura.beta-enabled", true)

	helper.ExecuteCommand("import job get --project-id=f607bebe-0cc0-4166-b60c-b4eed69ee7ee --job-id=87d485b4-73fc-4a7f-bb03-720f4672947e")

	mockHandler.AssertCalledTimes(1)
	mockHandler.AssertCalledWithMethod(http.MethodGet)

	helper.AssertErr("")
	helper.AssertOutJson(`
{
    "data": {
        "aura_target": {
            "db_id": "07e49cf5",
            "project_id": "f607bebe-0cc0-4166-b60c-b4eed69ee7ee"
        },
        "data_source": {
            "id": "177216b3-49b9-4e24-a414-a9ecbb448c53",
            "name": "AWS_POSTGRES_FLIGHTS",
            "type": "postgresql"
        },
        "id": "87d485b4-73fc-4a7f-bb03-720f4672947e",
        "import_type": "cloud",
        "info": {
            "cancellation_requested_time": "2025-08-15T13:15:19.655209Z",
            "completion_time": "2025-08-15T13:15:19Z",
            "exit_status": {
                "message": "Cancelled",
                "state": "Cancelled"
            },
            "last_update_time": "2025-08-18T12:25:53.469873Z",
            "percentage_complete": 95.23,
            "start_time": "2025-08-15T13:12:51Z",
            "state": "Completed",
            "submitted_time": "2025-08-15T13:12:50.499797Z"
        },
        "user_id": "f607bebe-0cc0-4166-b60c-b4eed69ee7ee"
    }
}
	`)
}

func TestGetImportJobByIdWithProgress(t *testing.T) {
	helper := testutils.NewAuraTestHelper(t)
	defer helper.Close()

	mockHandler := helper.NewRequestHandlerMock("/v2beta1/projects/f607bebe-0cc0-4166-b60c-b4eed69ee7ee/import/jobs/87d485b4-73fc-4a7f-bb03-720f4672947e", http.StatusOK, `
{
    "data": {
        "id": "87d485b4-73fc-4a7f-bb03-720f4672947e",
        "import_type": "cloud",
        "info": {
            "state": "Completed",
            "start_time": "2025-08-15T13:12:51Z",
            "completion_time": "2025-08-15T13:15:19Z",
            "exit_status": {
                "state": "Cancelled",
                "message": "Cancelled"
            },
            "cancellation_requested_time": "2025-08-15T13:15:19.655209Z",
            "submitted_time": "2025-08-15T13:12:50.499797Z",
            "last_update_time": "2025-08-18T12:25:53.469873Z",
            "percentage_complete": 95.23,
            "progress": {
                "nodes": [
                    {
                        "id": "n:1",
                        "labels": [
                            "AircraftData"
                        ],
                        "total_rows": 9,
                        "processed_rows": 9,
                        "created_nodes": 9,
                        "created_constraints": 1,
                        "created_indexes": 0
                    },
                    {
                        "id": "n:0",
                        "labels": [
                            "Aircraft"
                        ],
                        "total_rows": 9,
                        "processed_rows": 9,
                        "created_nodes": 9,
                        "created_constraints": 1,
                        "created_indexes": 0
                    },
                    {
                        "id": "n:3",
                        "labels": [
                            "AirportData"
                        ],
                        "total_rows": 104,
                        "processed_rows": 104,
                        "created_nodes": 104,
                        "created_constraints": 1,
                        "created_indexes": 0
                    },
                    {
                        "id": "n:2",
                        "labels": [
                            "Airport"
                        ],
                        "total_rows": 104,
                        "processed_rows": 104,
                        "created_nodes": 104,
                        "created_constraints": 1,
                        "created_indexes": 0
                    },
                    {
                        "id": "n:5",
                        "labels": [
                            "Flight"
                        ],
                        "total_rows": 33121,
                        "processed_rows": 33121,
                        "created_nodes": 33121,
                        "created_constraints": 1,
                        "created_indexes": 0
                    },
                    {
                        "id": "n:4",
                        "labels": [
                            "Booking"
                        ],
                        "total_rows": 262788,
                        "processed_rows": 262788,
                        "created_nodes": 262788,
                        "created_constraints": 1,
                        "created_indexes": 0
                    },
                    {
                        "id": "n:7",
                        "labels": [
                            "Ticket"
                        ],
                        "total_rows": 366733,
                        "processed_rows": 366733,
                        "created_nodes": 366733,
                        "created_constraints": 1,
                        "created_indexes": 0
                    },
                    {
                        "id": "n:6",
                        "labels": [
                            "FlightView"
                        ],
                        "total_rows": 33121,
                        "processed_rows": 33121,
                        "created_nodes": 33121,
                        "created_constraints": 1,
                        "created_indexes": 0
                    }
                ],
                "relationships": [
                    {
                        "id": "r:8",
                        "type": "TICKET_BELONGS_TO_BOOKING",
                        "total_rows": 366733,
                        "processed_rows": 255000,
                        "created_relationships": 255000,
                        "created_constraints": 0,
                        "created_indexes": 0
                    },
                    {
                        "id": "r:1",
                        "type": "FLIGHT_DEPARTS_FROM_AIRPORT",
                        "total_rows": 33121,
                        "processed_rows": 33121,
                        "created_relationships": 33121,
                        "created_constraints": 0,
                        "created_indexes": 0
                    },
                    {
                        "id": "r:0",
                        "type": "FLIGHT_USES_AIRCRAFT_DATA",
                        "total_rows": 33121,
                        "processed_rows": 33121,
                        "created_relationships": 33121,
                        "created_constraints": 0,
                        "created_indexes": 0
                    },
                    {
                        "id": "r:3",
                        "type": "FLIGHT_USES_AIRCRAFT",
                        "total_rows": 33121,
                        "processed_rows": 33121,
                        "created_relationships": 33121,
                        "created_constraints": 0,
                        "created_indexes": 0
                    },
                    {
                        "id": "r:2",
                        "type": "FLIGHT_ARRIVES_AT_AIRPORT",
                        "total_rows": 33121,
                        "processed_rows": 33121,
                        "created_relationships": 33121,
                        "created_constraints": 0,
                        "created_indexes": 0
                    },
                    {
                        "id": "r:5",
                        "type": "FLIGHT_VIEW_ARRIVES_AT_AIRPORT",
                        "total_rows": 33121,
                        "processed_rows": 33121,
                        "created_relationships": 33121,
                        "created_constraints": 0,
                        "created_indexes": 0
                    },
                    {
                        "id": "r:4",
                        "type": "FLIGHT_VIEW_DEPARTS_FROM_AIRPORT",
                        "total_rows": 33121,
                        "processed_rows": 33121,
                        "created_relationships": 33121,
                        "created_constraints": 0,
                        "created_indexes": 0
                    },
                    {
                        "id": "r:7",
                        "type": "FLIGHT_HAS_TICKET",
                        "total_rows": 1045726,
                        "processed_rows": 1045726,
                        "created_relationships": 1045726,
                        "created_constraints": 0,
                        "created_indexes": 0
                    },
                    {
                        "id": "r:6",
                        "type": "FLIGHT_VIEW_USES_AIRCRAFT",
                        "total_rows": 33121,
                        "processed_rows": 33121,
                        "created_relationships": 33121,
                        "created_constraints": 0,
                        "created_indexes": 0
                    }
                ]
            }
        },
        "data_source": {
            "id": "177216b3-49b9-4e24-a414-a9ecbb448c53",
            "type": "postgresql",
            "name": "AWS_POSTGRES_FLIGHTS"
        },
        "aura_target": {
            "db_id": "07e49cf5",
            "project_id": "f607bebe-0cc0-4166-b60c-b4eed69ee7ee"
        },
        "user_id": "f607bebe-0cc0-4166-b60c-b4eed69ee7ee"
    }
}
	`)

	helper.SetConfigValue("aura.beta-enabled", true)

	helper.ExecuteCommand("import job get --project-id=f607bebe-0cc0-4166-b60c-b4eed69ee7ee --job-id=87d485b4-73fc-4a7f-bb03-720f4672947e --progress=true")

	mockHandler.AssertCalledTimes(1)
	mockHandler.AssertCalledWithMethod(http.MethodGet)

	helper.AssertErr("")
	helper.AssertOutJson(`
{
    "data": {
        "aura_target": {
            "db_id": "07e49cf5",
            "project_id": "f607bebe-0cc0-4166-b60c-b4eed69ee7ee"
        },
        "data_source": {
            "id": "177216b3-49b9-4e24-a414-a9ecbb448c53",
            "name": "AWS_POSTGRES_FLIGHTS",
            "type": "postgresql"
        },
        "id": "87d485b4-73fc-4a7f-bb03-720f4672947e",
        "import_type": "cloud",
        "info": {
            "cancellation_requested_time": "2025-08-15T13:15:19.655209Z",
            "completion_time": "2025-08-15T13:15:19Z",
            "exit_status": {
                "message": "Cancelled",
                "state": "Cancelled"
            },
            "last_update_time": "2025-08-18T12:25:53.469873Z",
            "percentage_complete": 95.23,
            "progress": {
                "nodes": [
                    {
                        "created_constraints": 1,
                        "created_indexes": 0,
                        "created_nodes": 9,
                        "id": "n:1",
                        "labels": [
                            "AircraftData"
                        ],
                        "processed_rows": 9,
                        "total_rows": 9
					},
                    {
                        "created_constraints": 1,
                        "created_indexes": 0,
                        "created_nodes": 9,
                        "id": "n:0",
                        "labels": [
                            "Aircraft"
                        ],
                        "processed_rows": 9,
                        "total_rows": 9
                    },
                    {
                        "created_constraints": 1,
                        "created_indexes": 0,
                        "created_nodes": 104,
                        "id": "n:3",
                        "labels": [
                            "AirportData"
                        ],
                        "processed_rows": 104,
                        "total_rows": 104
                    },
                    {
                        "created_constraints": 1,
                        "created_indexes": 0,
                        "created_nodes": 104,
                        "id": "n:2",
                        "labels": [
                            "Airport"
                        ],
                        "processed_rows": 104,
                        "total_rows": 104
                    },
                    {
                        "created_constraints": 1,
                        "created_indexes": 0,
                        "created_nodes": 33121,
                        "id": "n:5",
                        "labels": [
                            "Flight"
                        ],
                        "processed_rows": 33121,
                        "total_rows": 33121
                    },
                    {
                        "created_constraints": 1,
                        "created_indexes": 0,
                        "created_nodes": 262788,
                        "id": "n:4",
                        "labels": [
                            "Booking"
                        ],
                        "processed_rows": 262788,
                        "total_rows": 262788
                    },
                    {
                        "created_constraints": 1,
                        "created_indexes": 0,
                        "created_nodes": 366733,
                        "id": "n:7",
                        "labels": [
                            "Ticket"
                        ],
                        "processed_rows": 366733,
                        "total_rows": 366733
                    },
                    {
                        "created_constraints": 1,
                        "created_indexes": 0,
                        "created_nodes": 33121,
                        "id": "n:6",
                        "labels": [
                            "FlightView"
                        ],
                        "processed_rows": 33121,
                        "total_rows": 33121
                    }
                ],
                "relationships": [
                    {
                        "created_constraints": 0,
                        "created_indexes": 0,
                        "created_relationships": 255000,
                        "id": "r:8",
                        "processed_rows": 255000,
                        "total_rows": 366733,
                        "type": "TICKET_BELONGS_TO_BOOKING"
                    },
                    {
                        "created_constraints": 0,
                        "created_indexes": 0,
                        "created_relationships": 33121,
                        "id": "r:1",
                        "processed_rows": 33121,
                        "total_rows": 33121,
                        "type": "FLIGHT_DEPARTS_FROM_AIRPORT"
                    },
                    {
                        "created_constraints": 0,
                        "created_indexes": 0,
                        "created_relationships": 33121,
                        "id": "r:0",
                        "processed_rows": 33121,
                        "total_rows": 33121,
                        "type": "FLIGHT_USES_AIRCRAFT_DATA"
                    },
                    {
                        "created_constraints": 0,
                        "created_indexes": 0,
                        "created_relationships": 33121,
                        "id": "r:3",
                        "processed_rows": 33121,
                        "total_rows": 33121,
                        "type": "FLIGHT_USES_AIRCRAFT"
                    },
                    {
                        "created_constraints": 0,
                        "created_indexes": 0,
                        "created_relationships": 33121,
                        "id": "r:2",
                        "processed_rows": 33121,
                        "total_rows": 33121,
                        "type": "FLIGHT_ARRIVES_AT_AIRPORT"
                    },
                    {
                        "created_constraints": 0,
                        "created_indexes": 0,
                        "created_relationships": 33121,
                        "id": "r:5",
                        "processed_rows": 33121,
                        "total_rows": 33121,
                        "type": "FLIGHT_VIEW_ARRIVES_AT_AIRPORT"
                    },
                    {
                        "created_constraints": 0,
                        "created_indexes": 0,
                        "created_relationships": 33121,
                        "id": "r:4",
                        "processed_rows": 33121,
                        "total_rows": 33121,
                        "type": "FLIGHT_VIEW_DEPARTS_FROM_AIRPORT"
                    },
                    {
                        "created_constraints": 0,
                        "created_indexes": 0,
                        "created_relationships": 1045726,
                        "id": "r:7",
                        "processed_rows": 1045726,
                        "total_rows": 1045726,
                        "type": "FLIGHT_HAS_TICKET"
                    },
                    {
                        "created_constraints": 0,
                        "created_indexes": 0,
                        "created_relationships": 33121,
                        "id": "r:6",
                        "processed_rows": 33121,
                        "total_rows": 33121,
                        "type": "FLIGHT_VIEW_USES_AIRCRAFT"
                    }
                ]
            },
            "start_time": "2025-08-15T13:12:51Z",
            "state": "Completed",
            "submitted_time": "2025-08-15T13:12:50.499797Z"
        },
        "user_id": "f607bebe-0cc0-4166-b60c-b4eed69ee7ee"
    }
}
	`)
}
