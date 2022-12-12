<p align="center">

<img  src="https://mk0abtastybwtpirqi5t.kinstacdn.com/wp-content/uploads/picture-solutions-persona-product-flagship.jpg"  width="211"  height="182"  alt="flagship-c-go-wrapper"  />

</p>

# Flagship Wrapper for C

## Description

Flagship wrapper for C expose functions compiled by Go SDK.

### Functions

#### init

- **syntax**: `init(char* environmentID, char* apiKey, int polling, char* logLevel, int trackingEnabled)`
- **default**: `none`
- **return**: `void`

Initialize the SDK with 5 arguments : Environment id, Api Key, Polling Interval, Log level and Tracking enabled.
Once the initialization is done the 1st time, the second you execute the script it will bypass the initialization.

#### getFlagBool

- **syntax**: `getFlagBool(char* visitorID, char* contextString, char* key, int defaultValue, int activate)`
- **default**: `none`
- **return**: `int`

Return flag that corresponds to the visitor id and context.

#### getFlagNumber

- **syntax**: `getFlagNumber(char* visitorID, char* contextString, char* key, double defaultValue, int activate)`
- **default**: `none`
- **return**: `double`

Return flag that corresponds to the visitor id and context.

#### getFlagString

- **syntax**: `getFlagString(char* visitorID, char* contextString, char* key, char* defaultValue, int activate)`
- **default**: `none`
- **return**: `char*`

Return flag that corresponds to the visitor id and context.

#### getAllFlags

- **syntax**: `getAllFlags(char* visitorID, char* contextString)`
- **default**: `none`
- **return**: `char*`

Return all flags that corresponds to the visitor id and context (flags are separated by ";").

### Dependencies

- Go

### Run

```bash
cd flagship-c-go-wrapper
chmod +x ./entrypoint.sh
./entrypoint.sh
``` 

Then copy the shared object generated in build folder.

## License

[Apache License.](https://github.com/flagship-io/flagship-c-go-wrapper/blob/main/LICENSE)