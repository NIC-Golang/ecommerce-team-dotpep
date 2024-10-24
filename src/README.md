# Services (Technical Requirements)

- User Auth service
- Core Shop service
- Scraper and AI service

## User Auth service

- User Authentication/Authorization with `JWT`, sending confirmation code via `SMTP` email.
- Password Recovery/Reset feature (Forgot Password) and it will send email to recovery/reset and set new password.
- Server Side cache storing for Auth in `Redis`, login expiration session date or time.
- Database `MongoDB` or `PostgreSQL`
- User Group/Role Management (It need to be Integrate/Communicate with `Core-Shop-Service`):
    - `Admin`, can assign Group/Role for other Clients/Users
    - `Manager`, can POST/Create or DELETE, and manipulate `Products`, `Category` Entities in `Core-Shop-Service`. Also assign a `Order` to a `Delivery` role.
    - `Delivery`, can handle Clients Order and check `ShippingAdress`, confirm that it was delivered.
    - `Client`/`Customer`/`Auth User`, can Add `Product` to `Bucket`/`OrderItem`, confirm that order has arrived.
    - `Unauth User`/`Not registered User`, can only View Product lists, Filter by Category them, and Register account as Client (Auth via JWT).

## Core Shop service

## Scraper and AI service
