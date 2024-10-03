# GoGent

A small application for building simple LLM based Agents using lua. GoGent is written in go and uses the [gopher-lua](https://github.com/yuin/gopher-lua) library to run lua code.

Currently GoGent is in a very early stage of development and is not yet ready for every day use. The goal is to provide a simple way to create LLM based agents that can be used in a variety of applications. Currently only [Ollama](https://github.com/ollama/ollama) is supported for all interactions.

## Roadmap

- [x] Basic lua support
- [x] Basic Ollama support
- [ ] Refactor Ollama support into LLM providers
- [ ] Basic provider support for Ollama, OpenAI and OpenAI through Azure
- [ ] Expose provider configuration
- [ ] Expose more functions to lua (memory, web requests and read/write files)
- [ ] Add support for [SQLite-vec](https://github.com/asg017/sqlite-vec)
