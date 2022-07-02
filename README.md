# html-tool
HTML-Tool criadas na linguagem Go by https://github.com/tomnomnom
# O que a ferramenta faz ?
Colocando URLs ou nomes de arquivo de documentos HTML em stdin, a ferramenta é usada como um filtro, ela procura na URLs/nome do arquivo todo os arquivos na página com o formato pre estabelecido! 
# Você vai precisar de GO para usar a ferramenta!
https://go.dev/dl/

apt install golang 
# Instalar a ferramenta!
```
 go get -u github.com/ELIZEUOPAIN/html-tool

```
# Como usar:
▶ html-tool 

```
Coloque  URLs ou nomes de arquivos para documentos HTML em stdin e extraia partes deles.

Usar: html-tool <comandos> [<args>]

comandos:
	tags <nome-da-tag>        Extrair texto contido em tags
	attribs <nome-dos-atributos>  Extrair valores dos atributos
	comments                Extrai comentários

Exemplo:
	echo "http://site.com.br" | /root/go/bin/html-tool comments 	*(Extrai os comentarios da pagina/URL)*
	echo "http://site.com.br" | /root/go/bin/html-tool tags title 	*(Extrai as tags Titulo da pagina/URL)*
	echo "http://site.com.br" | /root/go/bin/html-tool attribs src 	*(Extrai todos os atributos src da pagina/URL)*
_____________________________________________________________OU___________________________________________________________________________
	cat urls.txt | html-tool tags title a strong			*(Extrai as tags Titulo da pagina/URL)*
	find . -type f -name "*.html" | html-tool attribs src href 	*(Extrai todos os atributos src e href da pagina/URL)*
	cat urls.txt | html-tool comments 				*(extrai os comentarios da pagina/URL)*
```

OBS: As vezes não é possivel ultilizar a ferramenta ultilizando apenas o "html-tool" na sintaxe do comando no Linux, então ultilize o echo
| /root/go/bin/html-tool (diretório onde foi instalado a ferramenta) comando! 
# TODO

*Support selectors with https://github.com/ericchiang/css
        
*Option to ignore certificate errors
        
*Option to output filename/URL with output
