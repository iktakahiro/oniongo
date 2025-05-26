package di

import (
	"github.com/iktakahiro/oniongo/internal/application/todoapp"
	"github.com/iktakahiro/oniongo/internal/infrastructure/ent/repository/todorepo"
	"github.com/samber/do"
)

func DependencyInjection() *do.Injector {
	injector := do.New()

	// Repositories
	do.Provide(injector, todorepo.NewTodoRepository)

	// UseCases
	do.Provide(injector, todoapp.NewCreateTodoUseCase)
	do.Provide(injector, todoapp.NewGetTodoUseCase)
	do.Provide(injector, todoapp.NewGetTodosUseCase)
	do.Provide(injector, todoapp.NewUpdateTodoUseCase)
	do.Provide(injector, todoapp.NewStartTodoUseCase)
	do.Provide(injector, todoapp.NewCompleteTodoUseCase)
	do.Provide(injector, todoapp.NewDeleteTodoUseCase)

	return injector
}
