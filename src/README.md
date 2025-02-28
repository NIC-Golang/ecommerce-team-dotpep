# Services (Technical Requirements)

**In TODO:**

- User Auth service
- Core Shop service
- Scraper and AI service

---

**Common Requirements:**

- Stick and Follow Best Practices and Programming Principles.
- API Dev Best Practices:
    - Versioning (version you API)
    - Pagination (paginate lists of item endpoint get response)
    - OpenAPI/Swagger API Documentation
- Security:
    - Data Validation/Sanitization, Serialization/Deserialization.
    - Authentication and Authorization
    - Rate limiting (limits calls/requests, API throttling rate)
    - Middlewares
    - CORS
    - OWASP TOP 10 (SQL Injection prevention)
    - Firewall on your deployed server
- Performance:
    - DB Entity Indexing
    - Caching

## User Auth service

- Golang Gin for API
- JWT for User Auth

---

- User Authentication/Authorization with `JWT`, sending confirmation code via `SMTP` email.
- Password Recovery/Reset feature (Forgot Password) and it will send email to recovery/reset and set new password.
- Server Side cache storing for Auth in `Redis`, login expiration session date or time.
- Database `MongoDB` or `PostgreSQL`

---

**Postponed To Later (and for Consideration):**

- User Group/Role Management (It need to be Integrate/Communicate with `Core-Shop-Service`):
    - `Admin`, can assign Group/Role for other Clients/Users
    - `Manager`, can POST/Create or DELETE, and manipulate `Products`, `Category` Entities in `Core-Shop-Service`. Also assign a `Order` to a `Delivery` role.
    - `Delivery`, can handle Clients Order and check `ShippingAdress`, confirm that it was delivered.
    - `Client`/`Customer`/`Auth User`, can Add `Product` to `Bucket`/`OrderItem`, confirm that order has arrived.
    - `Unauth User`/`Not registered User`, can only View Product lists, Filter by Category them, and Register account as Client (Auth via JWT).
- Request, Rate limiting (limits calls/requests, API throttling rate) for users roles, endpoints (unauth user can request 5/minute, auth user 25/minute).

## Core Shop service

- Golang Gin for API
- PostgreSQL as Database (RDBMS)

---

- CRUD operations and API Endpoints:
    - Basic CRUD of E-Commerce Shop:
    - `one product` - get, create, delete, edit.
    - `list of products` - get.
    - `category` - create, delete.
    - `list of categories` - get.
    - `list of products in category` get `list of product` that have follow/relation with this category.
    - Unrequired for now (It might be separeted Service and in Consideration) (more CRUD in Management):
    - add `product` to `order item`.
    - get specific `user`, `order item` with `list of products` that was added.
    - process user `order item` and create `order`
    - `shipping adress` get, create, delete, edit.
    - create/add `product` to `favorite`.
    - get `list of products item` in `favorite`.
    - if user unauth then `order item` will saved in session cache using Redis or in Browser Cookies.
    - if user is authenticated then `order item` it will saved in Redis cache.
- Required PostgreSQL Tables/Entities/Model:
    - User (unrequired, temp) (but you can create it for testing and relating, it must be deleted later, when User-Auth-Service will be implemented and you will integrate it)
    - Product
    - Category (product category)
    - OrderItem (like bucket/basket or cart in shop, list of product that client picked)
    - Order (client order list of products that in OrderItem/Bucket/Basket/Cart)
    - ShippingAddress (unrequired for now)
    - FavoriteProduct (unrequired for now)
- Filtering, Ordering, Searching and Pagination

**Example:**

E-Commerce ER-D:

![E-Commerce ER-D](/docs/img/e-commerce_er-d.png)

---

Restaurant ER-D:

![Restaurant ER-D](/docs/img/restaurant_er-d.png)

## Scraper and AI service

- Golang Colly as Scraper
- OpenAI API for Filtering

---

**Scrap other service (Book seller Shops):**

- Scrape Amazon Shop, Flip.kz, Kaspi Shop, Meloman (marwin) kz websites (choose one of them first) Book section and save it in database Product table, using Golang Colly library, make http request and get html data, work with divs and html tags.

**Basic CRUD (REST API endpoints):**

- define what scraped data you need save, in your database table/entity.
- you can use MongoDB or PostgreSQL as Database.
- for REST API use Golang Gin
- get product by id
- get list of products
- Add Pagination feature
- Filter by price, descending and ascending order
