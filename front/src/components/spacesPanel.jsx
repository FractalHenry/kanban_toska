import { BoardLink } from "./links/boardlink"
import { UserLink } from "./links/userlink"
import { useNavigate } from "react-router-dom";
import { useState, useEffect } from "react";
import Cookies from 'js-cookie';
import { useToast } from "./Toast/toastprovider";

/* const spaces = [
  {
    name: "Space1",
    boards: [
      {
        id: 1,
        name: "Board1",
        color: "red"
      },
      {
        id: 2,
        name: "Board2",
        color: "blue"
      }, {
        id: 3,
        name: "Board3",
        color: "#EF23AA"
      }
    ],
    users: [
      {
        id: 1,
        login: "Star",
        username: "StarLord"
      },
      {
        id: 2,
        login: "CryBaby",
        username: "Drax"
      },
      {
        id: 3,
        login: "Chicka",
        username: "Gamora"
      }
    ]
  }, {
    name: "Space1",
    boards: [
      {
        id: 1,
        name: "Board1",
        color: "red"
      },
      {
        id: 2,
        name: "Board2",
        color: "blue"
      }, {
        id: 3,
        name: "Board3",
        color: "#EF23AA"
      }
    ],
    users: [
      {
        id: 1,
        login: "Star",
        username: "StarLord"
      },
      {
        id: 2,
        login: "CryBaby",
        username: "Drax"
      },
      {
        id: 3,
        login: "Chicka",
        username: "Gamora"
      }
    ]
  }
] */

const PanelItem = ({ space }) => {
  return (
    <div className="SpaceCard">
      <h4> Пространство: {space.name} </h4>
      <h5> Доски: </h5>
      {space.boards.map(
        (item, index) => (<BoardLink key={index} board={item} />)
      )}
      <h5>Участники пространства:</h5>
      {space.users.map(
        (item, index) => (<UserLink key={index} user={item} />)
      )}
    </div>
  )
}

export const SpacesPanel = () => {
  const navigate = useNavigate();
const { showToast } = useToast();

const { spaces, setSpaces } = useState();
useEffect(() => {
  const fetchData = async () => {
    const token = Cookies.get('authToken');
    if (!token) {
      navigate('/error/404');
      return;
    }
    try {
      const response = await fetch(`http://localhost:8000/spaces`, {
        method: 'GET',
        headers: {
          'Authorization': `Bearer ${token}`
        }
      });

      if (response.ok) {
        const data = await response.json();
        setSpaces(data);
      } else {
        throw new Error('Что-то пошло не так');
      }
    } catch (error) {
      showToast("Произошла ошибка при получении пространств");
    }
  };
  fetchData();
}, [navigate]);

  return (
    <div className="overflow-y w-min-300 h-max-100">
      <h3>Доступные пространства:</h3>
      {spaces&& spaces.map((item, index) => (<PanelItem key={index} space={item} />))}
    </div>
  )
}