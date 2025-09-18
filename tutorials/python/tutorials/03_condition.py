pwd_correct = "True"

if pwd_correct:
    print("---------------")
    print("Pass")
    print("---------------")
else:
    print("---------------")
    print("Fail")
    print("---------------")

## elseif ladder

if pwd_correct:
    print("True: Pass")
elif not pwd_correct:
    print("False: Fail")
else:
    print("Never: Pass")  ## Never prints if True since it executed if

num = int(input("Enter a number: "))
### Check multiple if statement
if 99 < num < 1000:
    print("Three Digit")
elif num > 999:
    print("More than Three Digit")
else:
    print("Less than Three Digit")

# Truth Table

name = "Ohman"

if name[0] == "R" or name[0] == "O":
    print("True")
elif name[0] == "R" and name[0] == "O":
    print("False")
else:
    print("Never: Pass")
### BitWise operator

## & , | , ~, ^, >> , <<
a = 21

# The >> operator moves all the bits in the number to the right.
# Each right shift divides the number by 2 (ignoring any remainder).
# For example, 21 in binary is 10101.
# Shifting right by 1 gives 1010, which is 10 in decimal.
b = a >> 1

print(b)  # This will print 10

# The code is optimized in these ways:
# - Bitwise operations (like >> and &) are very fast because they work directly on the binary representation of numbers.
# - Extracting RGB values using shifts and masks avoids slow string or math operations.
# - No unnecessary loops or conversions are used.
# - Each operation is done in a single line, making it both efficient and easy to read.

# Example: Extract RGB from 24-bit color integer
color = 0x12AC34  # Hexadecimal for Red=18, Green=172, Blue=52

# These bitwise operations are much faster than using division or modulo for extracting bytes.
red = (color >> 16) & 0xFF
green = (color >> 8) & 0xFF
blue = color & 0xFF

print(f"Red: {red}, Green: {green}, Blue: {blue}")
# This will print the amount of red, green, and blue in the color.
# To get the blue part, we just keep the last 8 bits.
blue = color & 0xFF

print(f"Red: {red}, Green: {green}, Blue: {blue}")
# This will print the amount of red, green, and blue in the color.

### Slice Strings

# %%
print("Hello")

# %% [markdown]
# # Lets Try this

# %%
print("Hello, JupyterLab!")


# %%
name = "Python"
## Slice syntax is string[start:stop:step]
# print(name[::-1])
print(name[3:1:-2])
# %%


## Lists
