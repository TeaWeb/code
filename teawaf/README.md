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

## Apply WAF
~~~
Request  -->  WAF  -->   Backends
			/
Response  <-- WAF <----		
~~~

## Coding
~~~go
waf := teawaf.NewWAF()

// add rule groups here

err := waf.Init()
if err != nil {
	return
}
waf.Start()

// stop the waf
// waf.Stop()
~~~