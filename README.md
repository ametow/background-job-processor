# Tasker – Job Scheduler  

### Project Overview  
This project processes scheduled tasks by making HTTP requests to third-party services. The client sends a task in JSON format, and the service returns a unique task ID. The task then runs in the background.  

### How It Works  
The system has two microservices:  
1. **Task Receiver:** Accepts tasks, stores them in a database, and provides task statuses.  
2. **Task Executor:** Regularly checks the database for new tasks, runs them in parallel (making HTTP requests), and updates their statuses.  

- [Flowchart](https://github.com/ametow/background-job-processor/blob/main/docs/diagram_microservices.pdf)  
- [Project Specification](https://github.com/ametow/background-job-processor/blob/main/docs/task_spec.pdf)  

### Configuration  
Set the PostgreSQL database connection string (DSN) in the environment variable **"DATABASE_DSN"**.  

### Running & Testing  

#### Run the Task Executor (Agent)  
- The agent runs on its own; no need to start the server.  
- Fill the database with test data:  
  - File: `internal/server/storage/db_test.go`  
  - Test name: `"TestTaskStorage_Create30Tasks_RealDB"`  
- Start the agent:  
  - Run `cmd/agent/agent.go`  

#### Run the Server  
- Start the server:
  - Prepare: set DATABASE_DSN environment variable
  - Run `cmd/main/main.go`  
- Test a **POST request**:  
  - File: `internal/server/storage/db_test.go`  
  - Test name: `"TestTaskStorage_CreateTask_RealDB"`  
- Test a **GET request**:  
  - File: `internal/server/storage/db_test.go`  
  - Test name: `"TestTaskStorage_GetTaskStatus_RealDB"`  

---