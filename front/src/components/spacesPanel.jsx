import { BoardLink } from "./links/boardlink"
import { UserLink } from "./links/userlink"

const spaces=[
  {
    name: "Space1",
    boards:[
      {
        id:1,
        name:"Board1",
        color: "red"
      },
      {
        id:2,
        name:"Board2",
        color:"blue"
      },{
        id:3,
        name:"Board3",
        color:"#EF23AA"
      }
    ],
    users:[
      {
        id:1,
        login: "Star",
        username:"StarLord"
      },
      {
        id:2,
        login:"CryBaby",
        username:"Drax"
      },
      {
        id:3,
        login: "Chicka",
        username:"Gamora"
      }
    ]
  },{
    name: "Space1",
    boards:[
      {
        id:1,
        name:"Board1",
        color: "red"
      },
      {
        id:2,
        name:"Board2",
        color:"blue"
      },{
        id:3,
        name:"Board3",
        color:"#EF23AA"
      }
    ],
    users:[
      {
        id:1,
        login: "Star",
        username:"StarLord"
      },
      {
        id:2,
        login:"CryBaby",
        username:"Drax"
      },
      {
        id:3,
        login: "Chicka",
        username:"Gamora"
      }
    ]
  }
]

const PanelItem = ({space}) =>{
  return(
    <div className="SpaceCard">
        <h4> Пространство: {space.name} </h4>
        <h5> Доски: </h5>
        {space.boards.map(
            (item,index)=>(<BoardLink key={index} board={item}/>)
        )}
        <h5>Участники пространства:</h5>
        {space.users.map(
            (item,index)=>(<UserLink key={index} user={item} />)
        )}
    </div>
    )
}


export const SpacesPanel = () =>{
    return (
        <div className="overflow-y w-min-300 h-max-100">
            <h3>Доступные пространства:</h3>
            {spaces.map((item,index)=>(<PanelItem key={index} space={item}/>))}
        </div>
    )
}