# make MAKE not to be confused by the frontend and server directory
.PHONY : server frontend

server:
	@echo starting the game server
	bash -c "cd server && go build . && ./server"

frontend: server
	@echo starting the game frontend
	bash -c "cd frontend && npm start"

# this is needed when you want a standalone frontend without starting the game server
frontend-only:
	@echo starting only the frontend
	bash -c "cd frontend && npm start"

run-game:
	make server & make frontend
