FROM mcr.microsoft.com/dotnet/sdk:10.0 AS build
WORKDIR /src

COPY DotnetSkaffold.Api.csproj ./
RUN --mount=type=cache,target=/root/.nuget/packages \
    dotnet restore DotnetSkaffold.Api.csproj

COPY . ./
RUN --mount=type=cache,target=/root/.nuget/packages \
    dotnet publish DotnetSkaffold.Api.csproj \
      --configuration Release \
      --no-restore \
      --output /app/publish

FROM mcr.microsoft.com/dotnet/aspnet:10.0-alpine AS runtime
WORKDIR /app

COPY --from=build /app/publish ./

USER $APP_UID
EXPOSE 8080
ENV ASPNETCORE_HTTP_PORTS=8080

ENTRYPOINT ["dotnet", "DotnetSkaffold.Api.dll"]
