# Ollama-openfunctions

This is a continuation of [ollama-functions](../ollama-functions/), using the [Gorilla OpenFunctions](https://gorilla.cs.berkeley.edu/blogs/4_open_functions.html)
model instead of a generic LLM. This seems to be much more stable in producing valid function calls. 

## Running

Create a `.env` file with the following contents:

```env
WEATHERAPI_KEY=<apikey>
```

(You can create an API-key by registering here: https://www.weatherapi.com/my/)

Run:

```bash
$ ./build.sh -b

$ ./bin/weather "is it raining in the capital of Italy?"
19:20:44 INFO Model: adrienbrault/gorilla-openfunctions-v2:Q4_K_M
19:20:44 INFO Query: is it raining in the capital of Italy?
19:20:54 INFO Call tool: getCurrentWeather("Rome, IT", "celsius")
19:21:06 INFO Final response: The current weather in Rome is partly cloudy with a temperature of 17°C.

$ ./bin/weather "what is the sum of the temperatures (in celcius) of Beijing and Walldorf?"
19:27:39 INFO Model: adrienbrault/gorilla-openfunctions-v2:Q4_K_M
19:27:39 INFO Query: what is the sum of the temperatures (in celcius) of Beijing and Walldorf?
19:27:40 INFO Call tool: getCurrentWeather("Beijing, China", "celsius")
19:27:41 INFO Call tool: getCurrentWeather("Walldorf, Germany", "celsius")
19:27:41 INFO Call tool: add(17, 12)
19:27:42 INFO Not a call: The sum of the temperatures (in Celsius) for Beijing and Walldorf is 29.
```

(It still doesn't work 100% of the time though)
