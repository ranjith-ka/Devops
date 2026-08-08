var builder = WebApplication.CreateBuilder(args);
var app = builder.Build();

app.MapGet("/", () => Results.Ok(new
{
    message = "Hello from .NET, Skaffold, and Kind!",
    environment = app.Environment.EnvironmentName,
    utcTime = DateTimeOffset.UtcNow
}));

app.MapGet("/healthz", () => Results.Ok(new { status = "healthy" }));

app.Run();
