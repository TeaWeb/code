# WAF
A basic WAF for TeaWeb.

## Config Constructions
~~~
WAF
  Inbound
	  Rule Groups
		Rule Sets
		  Rules
			Checkpoint Param <Operator> Value
  Outbound
  	  Rule Groups
  	    ... 				
~~~

# Apply WAF
~~~
Request  -->  WAF  -->   Backends
			/
Response  <-- WAF <----		
~~~