# Technical Requirements (E-Commerce Backend API)

`Техническое задание для студентов на Интернет Магазин Проект на Голанге`

There is 2 teams, for each 3 students and 1 mentor:

@dotpep, (github_link) Mentor - 3 students:
    - @, Максим (github_link), core service
    - @, Нуртай (github_link), auth/user service
    - @, Артем (github_link), scraper service

@alfinky, (github_link) Mentor - 3 stundents:
    - @, Жанболат (github_link), core service
    - @, Ергали (github_link), auth/user service
    - @, Бексултан (github_link), scraper service

---

Main goals of this project is to mentor.

## Table of Content

## TO DO for Mentors

- [ ] write Table of content
- [ ] additional Technical Requirements
- [ ] structure Boilerplate and Architectural dirs

## Principles and Rules

- Write Unit Test
- Write Documentation on `some-service/README.md`

Programming Principles:

- Clean Code and Clean Architecture
- SOLID
- DRY
- KISS
- YAGNI
- Naming Convention

## E-Commerce Microservices

Base service:

- Auth User service
- Core Shop service
- Scraper and AI service

Additional service:

- Payment system/service
- Notification system/service

Additional Components:

- Docker Containerization and Docker compose
- CI/CD Pipeline
- Centralized Logging system (ELK)

### Auth User service

### Core Shop service

### Scraper and AI service

## Tech Stack

- Golang
- Communcation between Services will be using RESTful API (in the future gRPC)
- `Gin` for API, `Colly` for Scraping, `OpenAI/Or other` for Filtering Scraped Product from other E-Shops.
- Communication with Databases (use whetever libs, methods you want, ORMs or Raw SQL query)
- PostgreSQL for `Core-Shop-Service`
- MongoDB for `User-Auth-Service`
- Cache Redis for `User-Authentication` session, `User-Order-Bucket` and etc.
- JWT for Authentication/Authorization

---

**Required:**

- RESTful API (required)
- Databases (required)
- Unit Testing (required)

---

**For Future:**

- gRPC and API Gateway (for future)
- Docker containers, Docker Compose and CI/CD Pipeline (for future)

## Some Explanation for Students

- You need to delete/cut `.example` name part in your files and make it own/use it.
- Some files like docker-compose.yml, and like api-gateway dir and etc. it's for future. First of all try to use and implement it in your local machine without using Docker and etc.
- You should work only on your provided service `some-service/` directory.

---

**Life-cicle:**

`Deadline is 1 month, start date: 26.10.2024, end date: 30.11.2024`

- First step: use `git clone`
- You will just make pull request, with changes in your `some-service/`, push it.
- And We Mentors will be `code review` your code and apply or decline.
- Also 1 week 1 meetup, We will make some lections and activities, you will explain your written code, and etc. communication things.
- Students have TODO Tasks in Project Kanban Board assigned from mentor. It will be like 1 task 1 week, with issued in project and assigned to student.
- Mentor cannot write any code!

**ER-D Examples (for Core-Shop-Service):**

E-Commerce ER-D:

![E-Commerce ER-D](docs/assets/e-commerce_er-d.png)

---

Restaurant ER-D:

![Restaurant ER-D](docs/assets/restaurant_er-d.png)

**Links:**

- https://roadmap.sh/projects/ecommerce-api
- https://roadmap.sh/projects/scalable-ecommerce-platform

---

- https://drawsql.app/teams/test-team1-2/diagrams/littlelemonapi
- https://drawsql.app/teams/test-team1-2/diagrams/ecommerce
