# WAF
A basic WAF for TeaWeb.

## Config Constructions
~~~
WAF
  Action Configs
  Rule Groups
    Rule Sets
      Rules
        Check Point
~~~

# Apply WAF
~~~
Request  -->  WAF  -->   Backends
						/
Response  <-- WAF <----		
~~~