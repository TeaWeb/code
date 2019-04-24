# WAF
A basic WAF for TeaWeb.

## Config Constructions
~~~
WAF
  Rule Groups
    Rule Sets
      Rules
        Checkpoint Param <Operator> Value
~~~

# Apply WAF
~~~
Request  -->  WAF  -->   Backends
			/
Response  <-- WAF <----		
~~~