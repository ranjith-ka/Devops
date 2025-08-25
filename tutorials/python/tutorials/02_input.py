## To get input from user

name = input("What is your name? ")
print("Hello, " + name)

## Note: Print as like discussions

height = input("What is your height? ")
print(type(height)) # Default sting if input
print("Your height is " + height + " cm")

## Get input as float
height = float(input("What is your height again? "))
height_inches = "{:.2f}".format(height/2.54)
print(height_inches)
print("----------------------")
print("Your height is " + str(height) + " cm")
print("----------------------")
print("Your height is " + str(height_inches) + " cm")
print("----------------------")
print(round(float(height_inches)))

# Exercise
### Name, email, Phone
# TODO

## Maths/Execution

# It's not a compiled program, it executes over line by line

## (first_do_this)+next

## Maths functions
# round(), abs(a), pow(a,3)

import math
print("------------math----------")
print(math.ceil(float(height_inches)))
