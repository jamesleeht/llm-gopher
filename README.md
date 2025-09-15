# LLM Gopher

This is a small library that acts as a router and adapter for different LLM endpoints.

It doesn't provide any big framework capabilities - it is simply meant as an easy way to create a router for different models.

It is partially inspired by the (LiteLLM Python SDK)[https://github.com/BerriAI/litellm] but in a Golang context.

## Client types supported

- OpenAI
- Vertex

## Presets

A preset represents a combination of the model and its settings.
We use this internally since we have different use cases which might require a specific combination.

## Router

The router can be configured with a bunch of clients and presets.

1. When you send a prompt to the router, you specify a preset.
2. The preset's settings will be applied and an appropriate client will be selected for the model.
