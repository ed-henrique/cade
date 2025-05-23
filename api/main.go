package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	_ "github.com/joho/godotenv/autoload"
)

const (
	urlBase                 = "https://cade.ed-henrique.com"
	apiCorreiosAutenticacao = "https://api.correios.com.br/token/v1/autentica"
	apiCorreiosRastreamento = "https://apps3.correios.com.br/areletronico/v1/ars/eventos"
)

var (
	usuarioCorreios = os.Getenv("CORREIOS_USERNAME")
	senhaCorreios   = os.Getenv("CORREIOS_PASSWORD")
	credenciais     = base64.StdEncoding.EncodeToString([]byte(usuarioCorreios + ":" + senhaCorreios))
)

type TimeWithoutTimezone time.Time

type evento struct {
	TipoEvento              string `json:"tipoEvento"`
	StatusEvento            string `json:"statusEvento"`
	DescricaoEvento         string `json:"descricaoEvento"`
	NomeUnidade             string `json:"nomeUnidade"`
	Municipio               string `json:"municipio"`
	Uf                      string `json:"uf"`
	DataCriacao             string `json:"dataCriacao"`
	Latitude                string `json:"latitude"`
	Longitude               string `json:"longitude"`
	NomeRemetente           string `json:"nomeRemetente"`
	CepRemetente            string `json:"cepRemetente"`
	LogradouroRemetente     string `json:"logradouroRemetente"`
	ComplementoRemetente    string `json:"complementoRemetente"`
	NumeroRemetente         string `json:"numeroRemetente"`
	BairroRemetente         string `json:"bairroRemetente"`
	CidadeRemetente         string `json:"cidadeRemetente"`
	UfRemetente             string `json:"ufRemetente"`
	PaisRemetente           string `json:"paisRemetente"`
	NomeDestinatario        string `json:"nomeDestinatario"`
	CepDestinatario         string `json:"cepDestinatario"`
	LogradouroDestinatario  string `json:"logradouroDestinatario"`
	ComplementoDestinatario string `json:"complementoDestinatario"`
	EmailDestinatario       string `json:"emailDestinatario"`
	NumeroDestinatario      string `json:"numeroDestinatario"`
	BairroDestinatario      string `json:"bairroDestinatario"`
	CidadeDestinatario      string `json:"cidadeDestinatario"`
	UfDestinatario          string `json:"ufDestinatario"`
	PaisDestinatario        string `json:"paisDestinatario"`
	NomeRecebedor           string `json:"nomeRecebedor"`
	DataRecebimento         string `json:"dataRecebimento"`
	Documento               string `json:"documento"`
	Matricula               string `json:"matricula"`
	Usuario                 string `json:"usuario"`
	CodigoSRO               string `json:"codigoSRO"`
}

type objeto struct {
	Codigo             string   `json:"codigo"`
	ImagemBase64       string   `json:"imagemBase64"`
	Mensagem           string   `json:"mensagem"`
	Eventos            []evento `json:"eventos"`
	Tipo               string   `json:"tipo"`
	TipoEventoImagem   string   `json:"tipoEventoImagem"`
	DataCriacaoImagem  string   `json:"dataCriacaoImagem"`
	StatusEventoImagem string   `json:"statusEventoImagem"`
}

type acesso struct {
	Token    string              `json:"token"`
	ExpiraEm TimeWithoutTimezone `json:"expiraEm"`
}

func (t *TimeWithoutTimezone) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	tt, err := time.Parse("2006-01-02T15:04:05", s)
	if err != nil {
		return err
	}

	*t = TimeWithoutTimezone(tt)
	return nil
}

func (t TimeWithoutTimezone) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Time(t))
}

func (t TimeWithoutTimezone) Format(s string) string {
	tt := time.Time(t)
	return tt.Format(s)
}

func handleErr(w http.ResponseWriter, serverMsg, clientMsg string, code int, err error) {
	w.Header().Set("content-type", "text/plain")
	slog.Error(serverMsg, slog.String("err", err.Error()))
	http.Error(w, clientMsg, code)
}

func HandleServerErr(w http.ResponseWriter, serverMsg string, code int, err error) {
	handleErr(w, serverMsg, "Houve um erro no servidor.", code, err)
}

func createAndExecuteRequest(endpoint string, w http.ResponseWriter, r io.Reader, obj any) error {
	req, err := http.NewRequest(http.MethodPost, endpoint, r)
	if err != nil {
		HandleServerErr(w, "could not create new request to Correios", http.StatusInternalServerError, err)
		return err
	}

	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("authorization", "Basic "+credenciais)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		HandleServerErr(w, "could not generate access token", http.StatusInternalServerError, err)
		return err
	}
	defer res.Body.Close()

	if err := json.NewDecoder(res.Body).Decode(obj); err != nil {
		handleErr(w, "could not decode json", "Algum dos códigos enviados estava errado.", http.StatusBadRequest, err)
		return err
	}

	return nil
}

func main() {
	var acessoAtual acesso

	http.HandleFunc("OPTIONS /rastreamento", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("access-control-allow-origin", urlBase)
		w.Header().Add("access-control-allow-methods", strings.Join([]string{http.MethodPost, http.MethodOptions}, ", "))
		w.Header().Add("access-control-allow-headers", "Authorization, Accept, Content-Type")
	})

	http.HandleFunc("POST /rastreamento", func(w http.ResponseWriter, r *http.Request) {
		objetosRastreamento := struct {
			Objetos []string `json:"objetos"`
		}{}

		w.Header().Add("content-type", "application/json")
		w.Header().Add("access-control-allow-origin", urlBase)
		w.Header().Add("access-control-allow-methods", strings.Join([]string{http.MethodPost, http.MethodOptions}, ", "))
		w.Header().Add("access-control-allow-headers", "Authorization, Accept, Content-Type")
		w.Header().Add("vary", "Origin")

		defer r.Body.Close()
		if err := json.NewDecoder(r.Body).Decode(&objetosRastreamento); err != nil {
			handleErr(w, "could not decode json", "Algum dos códigos enviados estava errado.", http.StatusBadRequest, err)
			return
		}

		if time.Now().After(time.Time(acessoAtual.ExpiraEm)) {
			if err := createAndExecuteRequest(apiCorreiosAutenticacao, w, nil, &acessoAtual); err != nil {
				return
			}
		}

		var buffer bytes.Buffer
		if err := json.NewEncoder(&buffer).Encode(objetosRastreamento); err != nil {
			HandleServerErr(w, "could not encode json", http.StatusInternalServerError, err)
			return
		}

		objs := make([]objeto, len(objetosRastreamento.Objetos))
		if err := createAndExecuteRequest(apiCorreiosRastreamento, w, nil, &objs); err != nil {
			return
		}

		if _, err := io.Copy(w, &buffer); err != nil {
			HandleServerErr(w, "could not send full response", http.StatusInternalServerError, err)
			return
		}
	})

	if err := http.ListenAndServe(":8080", nil); err != nil {
		slog.Error("could not start server", slog.String("err", err.Error()))
		os.Exit(1)
	}
}
