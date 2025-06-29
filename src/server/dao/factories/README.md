# Factory Pattern for Testing

This directory contains factory functions to create model instances for testing purposes. The factories follow a consistent pattern to make test data creation easier and more maintainable.

## Available Factories

### Core Entities
- `user.go`: Factory for creating User models
- `community.go`: Factory for creating Community models
- `plan.go`: Factory for creating Plan models
- `membership.go`: Factory for creating Membership models
- `onboarding.go`: Factory for creating Onboarding models

### Services
- `service.go`: Factory for creating Service models
- `service_local.go`: Factory for creating ServiceLocal models
- `service_professional.go`: Factory for creating ServiceProfessional models
- `local.go`: Factory for creating Local models
- `professional.go`: Factory for creating Professional models
- `template.go`: Factory for creating Template models

### Reservations
- `session.go`: Factory for creating Session models
- `reservation.go`: Factory for creating Reservation models

### Relations
- `community_plan.go`: Factory for creating CommunityPlan models
- `community_service.go`: Factory for creating CommunityService models

## How to Use

Each factory provides two main functions:

1. `NewXXXModel`: Creates a single instance of the model
2. `NewXXXModelBatch`: Creates multiple instances of the model

### Basic Usage

```go
// Create a user with default values
user := factories.NewUserModel(db)

// Create a user with custom values
name := "Custom Name"
email := "custom@example.com"
rol := model.UserRolAdmin

customUser := factories.NewUserModel(db, factories.UserModelF{
    Name:  &name,
    Email: &email,
    Rol:   &rol,
})

// Create multiple users
users := factories.NewUserModelBatch(db, 5)
```

### Creating Related Entities

Some factories automatically create related entities. For example, the `NewMembershipModel` factory creates a User, Community, and Plan if not provided:

```go
// Create a membership with default related entities
membership := factories.NewMembershipModel(db)

// Create a membership with a specific user
userId := existingUser.Id
membership := factories.NewMembershipModel(db, factories.MembershipModelF{
    UserId: &userId,
})
```

## Best Practices

1. Use factories instead of manually creating model instances in tests
2. Create variables for string values before passing them to factory options
3. For enum types like `UserRol`, create a variable and pass its address
4. Use batch creation when you need multiple similar instances
5. When testing specific scenarios, override only the necessary fields and let the factory handle the rest 