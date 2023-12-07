import { createApp } from 'vue';
import "./assets/styles/_datatable.scss";
import PrimeVue from "primevue/config";
import 'primevue/resources/primevue.min.css';
import "primeflex/primeflex.css";
import 'remixicon/fonts/remixicon.css';
import 'primeicons/primeicons.css';

import 'primevue/resources/themes/lara-dark-cyan/theme.css';
import 'primevue/resources/themes/lara-light-cyan/theme.css';

import router from "./router";
import App from "./App.vue";
import Button from "primevue/button";
import Calendar  from "primevue/calendar";
import Checkbox  from "primevue/checkbox";
import Column  from "primevue/column";
import DataView from "primevue/dataview";
import Datatable from "primevue/datatable";
import Dialog from "primevue/dialog";
import Divider from "primevue/divider";
import InputText from "primevue/inputtext";
import TabPanel from 'primevue/tabpanel';
import TabView from 'primevue/tabview';
import Tooltip from "primevue/tooltip";
import Paginator from "primevue/paginator";
import Panel from "primevue/panel";
import ProgressBar from 'primevue/progressbar';
import RadioButton from 'primevue/radiobutton';
import StyleClass from "primevue/styleclass";
import Splitter from "primevue/splitter";
import SplitterPanel from "primevue/splitterpanel";
import Tag from "primevue/tag";
import TextArea from "primevue/textarea";
import Toolbar from "primevue/toolbar";

const app = createApp(App);
app.use(PrimeVue, {
  locale: {
    startsWith: "Commence par",
    contains: "Contient",
    notContains: "Ne contient pas",
    endsWith: "Se termine par",
    equals: "Égal",
    notEquals: "Non égal",
    noFilter: "Pas de filtre",
    lt: "Inférieur à",
    lte: "Inférieur ou égal à",
    gt: "Supérieur à",
    gte: "Supérieur ou égal à",
    dateIs: "La date est",
    dateIsNot: "La date n'est pas",
    dateBefore: "La date est avant",
    dateAfter: "La date est après",
    clear: "Effacer",
    apply: "Appliquer",
    matchAll: "Tout correspond",
    matchAny: "N'importe quel correspond",
    accept: "Accepter",
    reject: "Rejeter",
    choose: "Choisir",
    cancel: "Annuler",
    completed: "Terminé",
    pending: "En attente",
    fileSizeTypes: ["O", "Ko", "Mo", "Go", "To", "Po", "Eo", "Zo", "Yo"],
    dayNames: [
      "Dimanche",
      "Lundi",
      "Mardi",
      "Mercredi",
      "Jeudi",
      "Vendredi",
      "Samedi",
    ],
    dayNamesShort: ["Dim", "Lun", "Mar", "Mer", "Jeu", "Ven", "Sam"],
    dayNamesMin: ["Di", "Lu", "Ma", "Me", "Je", "Ve", "Sa"],
    monthNames: [
      "Janvier",
      "Février",
      "Mars",
      "Avril",
      "Mai",
      "Juin",
      "Juillet",
      "Août",
      "Septembre",
      "Octobre",
      "Novembre",
      "Décembre",
    ],
    monthNamesShort: [
      "Jan",
      "Fév",
      "Mar",
      "Avr",
      "Mai",
      "Juin",
      "Juil",
      "Août",
      "Sep",
      "Oct",
      "Nov",
      "Déc",
    ],
    chooseYear: "Choisir l'année",
    chooseMonth: "Choisir le mois",
    chooseDate: "Choisir la date",
    prevDecade: "Décennie précédente",
    nextDecade: "Décennie suivante",
    prevYear: "Année précédente",
    nextYear: "Année suivante",
    prevMonth: "Mois précédent",
    nextMonth: "Mois suivant",
    am: "AM",
    pm: "PM",
    today: "Aujourd'hui",
    weekHeader: "Sm",
    firstDayOfWeek: 1,
    showMonthAfterYear: false,
    dateFormat: "dd/mm/yy",
    weak: "Faible",
    medium: "Moyen",
    strong: "Fort",
    passwordPrompt: "Saisissez un mot de passe",
    emptyFilterMessage: "Aucun résultat trouvé",
    searchMessage: "{0} résultats sont disponibles",
    selectionMessage: "{0} éléments sélectionnés",
    emptySelectionMessage: "Aucun élément sélectionné",
    emptySearchMessage: "Aucun résultat trouvé",
    emptyMessage: "Aucune option disponible",
  },
  theme: 'lara-light-cyan'
});

app.use(router);

app.directive("tooltip", Tooltip);
app.directive("styleclass", StyleClass);

app.component("MyButton", Button);
app.component("MyCalendar", Calendar);
app.component("MyCheckbox", Checkbox);
app.component("MyColumn", Column);
app.component("DataView", DataView);
app.component("MyDatatable", Datatable);
app.component("MyDivider", Divider);
app.component("InputText", InputText);
app.component("MyDialog", Dialog);
app.component("MyPaginator", Paginator);
app.component("MyPanel", Panel);
app.component("MyProgressBar", ProgressBar)
app.component("MyRadioButton", RadioButton);
app.component("MyTabPanel", TabPanel);
app.component("MyTabView", TabView);
app.component("MyToolbar", Toolbar);
app.component("MySplitter", Splitter);
app.component("MySplitterPanel", SplitterPanel);
app.component("MyTag", Tag);
app.component("TextArea", TextArea);

app.mount("#app");