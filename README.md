# Image Gallery

Web Application that hosts a simple image gallery on your computer. Built with Vue.js as frontend with bootstrap and Sass to handle CSS. Golang for backend with a SQL database to store data.

# Running and compiling

## Setting up the database

- Start your a SQL environment
- Run the SchemaSetup.sql file to setup the database

## Starting the backend

- Navigate to the backend folder ```cd ./backend```
- Run ```go run . -- DBUSER DBPASS```, where ```DBUSER``` is the username for your SQL environment and ```DBPASS``` is the password, to start the backend server.

## Starting the frontend

- Navigate to the frontend folder ```cd ./frontend```
- Run ```npm install``` to add the necessary libraries
- Run ```npm run build``` to start the frontend server