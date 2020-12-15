# golang html template examples

* layout
* block
* include
* control
* loop
* actions
* pipelines
* functions



{{ define "content" }}  
<h1>You are awesome  {{ .Title }}</h1>  

{{ range $key,$value := .Datas }}  
    
    {{ if gt $value 10 }}  
    <li>{{ $value }} **** {{ $key }}</li>  
    {{ end }}  
   
{{ end }}  

<h3>  
Pieline: first func return result  is the next func argument   
</h3>  
{{ "One func" | printf "This is 2  %q" | printf "This is 3 %q"}}  
<br>
{{ with $x := "output" | printf "%q" }}{{ $x }}{{ end }}  
	A with action that creates and uses a variable.  

<h3>And</h3>  
{{ and 0 1 }}  

<h3>call</h3>  
{{ call .Mul 9 9 }}  

<h3>html index js lan </h3>  
{{ html "<h2> You are awesome </h2>" }} <br>
{{ index .Datas }} <br>
{{ js "<script>alert('You are awesome')</script>" }} <br>
{{ len "123456789哈哈" }} <br>

<h3>and not or , eq ne lt le gt ge</h3>  
{{ and 1 2  }} <br>
{{ and false 2  }}<br>
{{ not false  }}<br>
{{ or 0 1  }}<br>

{{ if gt 3 6 }}  
3>6  
{{ else }}  
3  less than 6  
{{ end }}  


<br>
{{ end }}
