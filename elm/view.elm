import Set
import List(map)
import List
import Signal

import Color 
import Graphics.Collage as Collage
import Graphics.Element(Element)
import Graphics.Collage(Shape, Form)
import Graphics.Element as Element
import Window
import Touch
import Text

import Abalone 
import Hex
import Player
import Player(Player(..))
import Misc

import Debug

type alias WidthHeight = (Int, Int)
type alias HexSize = Float

radiusInHex = 5

main : Signal Element
main = Signal.map2 scene Window.dimensions Touch.taps

scene : (Int, Int) -> {x:Int, y:Int } -> Element
scene (w,h) {x,y} = 
    let backgroundBoard = board (w,h) game
        game = Abalone.start
        marbles = stones (w,h) game 
    --positioned = move (toFloat x - toFloat w/2, toFloat h/2 - toFloat y)
    --    taps = collage w h [positioned (filled purple (circle 40)) ]
    in  Element.layers [backgroundBoard, marbles]

margin = 0

hexagon : HexSize -> Shape
hexagon = Collage.ngon 6

reposition : HexSize -> Hex.Position -> Form -> Form
reposition pixelRadius (q, r) hex = 
    let qf = toFloat q
        rf = toFloat r
        x = pixelRadius * sqrt(3) * (qf + rf/2)
        y = pixelRadius * 3/2 * rf
    in  Collage.move (x, y) hex

hexSize : WidthHeight -> Abalone.Game -> HexSize
hexSize (w, h) game = let wholeBoardLen = toFloat (min w h - margin)
                          hexesOnEdge = game.board.boardRadius
          in wholeBoardLen / (toFloat <| hexesOnEdge * 4) -- discovered experimentally ;)

genHex : HexSize -> Hex.Position -> Form
genHex size pos = 
    let style = Collage.outlined <| Collage.solid Color.black
        hex = Collage.rotate (degrees 30) <| style <| hexagon size 
        coord = Collage.toForm <| Text.plainText <| toString pos
        gp : Form 
        gp = Collage.group [hex, coord]
    in  reposition size pos gp

stone : HexSize -> Player -> Hex.Position -> Form
stone size p pos = Collage.circle (size / 2) |> Collage.filled (Player.colorOf p) |> reposition size pos

stones : WidthHeight -> Abalone.Game -> Element
stones (w, h) game = 
    let size = hexSize (w, h) game
        f = (\player -> map (stone size player) <| Set.toList <| Abalone.getPieces game.board player)
    in  Collage.collage w h <| List.concat <| map f [White, Black]

board : WidthHeight -> Abalone.Game -> Element
board (w, h) game = 
    let hexesOnEdge = game.board.boardRadius
        size = hexSize (w, h) game
        hexPositions : List Hex.Position
        hexPositions = Hex.hexagonalGrid hexesOnEdge
        hexagons : List Form
        hexagons = map (genHex size) hexPositions
    in  Collage.collage w h hexagons