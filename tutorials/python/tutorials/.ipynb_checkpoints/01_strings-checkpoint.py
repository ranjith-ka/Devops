message = "new projects"
name = "rk"

print("message", message)
print("name", name)
print(name)
print(name.upper())
print(message.lower())
print(message.title())

## Length of messages
print(len(message))
print(message.capitalize())

## -1 will return if character not available
print(message.find("k"))
print(message.count("k"))
print(message.isdigit())

print("Multiple assignments")
print("----------------------")
## Multiple assignments
name, age, fgr = "Rk", 30, "tt"
print(name)


print("Type casting")
print("----------------------")
## Type casting

## print("Age" + age) -Temp casting
print("Age: " + str(age)) # Type Casting (Convert int to sting for a temp)
print(5 + age)
