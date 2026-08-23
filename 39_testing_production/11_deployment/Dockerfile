# ---------- STAGE 1: BUILD ----------
# a full Go toolchain image, used only to compile the binary - none of this
# stage's ~800MB of compiler/tooling ends up in the final image below
FROM golang:latest AS builder

WORKDIR /src

# copying go.mod (and go.sum, if there were external dependencies) BEFORE
# the rest of the source lets Docker cache the module download step - it
# only re-runs when go.mod/go.sum actually change, not on every code edit
COPY go.mod ./
# no `go mod download` needed here - this module has zero external
# dependencies, just the standard library. a real project with dependencies
# would add: RUN go mod download

COPY 11_deployment/ ./11_deployment/

# CGO_ENABLED=0 produces a fully static binary with no dynamic libc
# dependency - the whole reason this can run on a minimal/scratch base at
# all without a "wrong libc version" failure
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/server ./11_deployment

# ---------- STAGE 2: RUN ----------
# alpine instead of the ~800MB builder image above - just enough OS to run
# a static binary and give HEALTHCHECK a shell/wget to use. `FROM scratch`
# would be even smaller (a few MB, just the binary) but has no shell at all,
# which makes HEALTHCHECK/debugging harder - a real tradeoff, not a strict upgrade
FROM alpine:latest

# a container running as root is a bigger blast radius if the app is ever
# compromised - create and switch to an unprivileged user instead
RUN adduser -D -H appuser
USER appuser

COPY --from=builder /out/server /server

ENV PORT=8080
EXPOSE 8080

# lets `docker ps` and orchestrators see "healthy"/"unhealthy" directly,
# using the /healthz endpoint this app exposes for exactly this purpose
HEALTHCHECK --interval=10s --timeout=3s \
  CMD wget -q -O- http://localhost:8080/healthz || exit 1

ENTRYPOINT ["/server"]
