# Lorenzo Codes API
This project is a RESTful API built with Go, designed to serve as the backend for a hosted frontend resume application. It provides endpoints to manage and retrieve data for Projects, Work History, and Certifications, stored in AWS DynamoDB. The API is intended to be deployed (e.g., on AWS Elastic Beanstalk) and consumed by a frontend (e.g., a React app) to display a dynamic, data-driven resume.

## Features
CRUD Operations: Create, read, update, and delete entries for Projects, Work History, and Certifications.
AWS DynamoDB: Persistent storage for resume data.
CORS Support: Allows requests from a frontend hosted on a different origin (e.g., http://localhost:3000).