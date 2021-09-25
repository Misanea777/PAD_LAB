using System;
using System.Collections.Generic;
using UnityEngine;

public class Instantiater : MonoBehaviour
{
    public int gridHeight;
    private int gridWidth;
    private float tileSize;
    Tile [ , ] tiles;
    public GameObject templateGreenTile;
    public GameObject templateVioletTile;


    HashSet<Chunck> chuncks = new HashSet<Chunck>();
    // Start is called before the first frame update
    void Start()
    {
        gridWidth = Mathf.RoundToInt(gridHeight*Camera.main.aspect);
        tileSize = (Camera.main.orthographicSize * 2) / gridHeight;

        
        chuncks.Add(Test.genChunck(new Vector2(0,0)));
        
        RenderChuncks(chuncks);

    }

    // Update is called once per frame
    void Update()
    {
        CheckIfChuncksVisible(chuncks);
        DestroyIfNotVisible(chuncks);
        Controller(chuncks);
    }

    void CheckIfChuncksVisible(HashSet<Chunck> chuncks) {
        Camera camera = Camera.main;
        Vector3 cameraPoint1 = camera.ViewportToWorldPoint(new Vector3(0,1,camera.nearClipPlane));
        Vector3 cameraPoint2 = camera.ViewportToWorldPoint(new Vector3(1,0,camera.nearClipPlane));
        foreach(Chunck c in chuncks) {
            c.CheckIfVisible(cameraPoint1, cameraPoint2, tileSize);
        }
    }

    void DestroyIfNotVisible(HashSet<Chunck> chuncks) {
        foreach(Chunck c in chuncks) {
            if(!c.isVisible && c.isRendered) {
                foreach(GameObject tile in c.renderedTiles) Destroy(tile);
                c.isRendered = false;
            }
        }
    }

    void RenderChuncks(HashSet<Chunck> chuncks) {
        foreach(Chunck c in chuncks) {
            if(c.isRendered) {
                continue;
            }
            RenderChunck(c);
        }
    }

    void RenderChunck(Chunck c) {
        if(c.isRendered) return;

        List<GameObject> renderedTiles = new List<GameObject>();
        foreach(Tile tile in c.floorTiles) {
            Vector3 tilePos = new Vector3(
                tile.posY * tileSize + tileSize/2,
                (tileSize * gridHeight) - (tile.posX * tileSize + tileSize/2),
                0 
            );
            renderedTiles.Add(Instantiate(TemplateFor(tile.type), tilePos, Quaternion.identity));
        }
        c.renderedTiles = renderedTiles;
        c.isRendered = true;
    }

    GameObject TemplateFor(TileType type) {
        if(type.Value == TileType.Green.Value) return templateGreenTile;
        if(type.Value == TileType.Violet.Value) return templateVioletTile;
        
        return templateGreenTile;
    }

    void Controller(HashSet<Chunck> chuncks) {
        Vector3 dir = new Vector3(0,0,0);
        int speed = 10;
        if(Input.GetKey(KeyCode.RightArrow))
        {
            dir = new Vector3(-speed * Time.deltaTime,0,0);
        }
        if(Input.GetKey(KeyCode.LeftArrow))
        {
            dir = new Vector3(speed * Time.deltaTime,0,0);
        }
        if(Input.GetKey(KeyCode.DownArrow))
        {
            dir = new Vector3(0,speed * Time.deltaTime,0);
        }
        if(Input.GetKey(KeyCode.UpArrow))
        {
           dir =  new Vector3(0,-speed * Time.deltaTime,0);
        }


        foreach(Chunck chunck in chuncks) {
            Chunck c = chunck;

            if(!c.isRendered) {
                foreach(GameObject tile in c.renderedTiles) Destroy(tile);
                continue;
            }
            c.translate(dir, tileSize);
            List<Chunck> nextToRender =  c.getNextChuncks();
            // foreach(Chunck nextC in nextToRender) chuncks.Add(nextC);
        }
    }

}

abstract class Tile
{

    public int posX {get; set;}  
    public int posY {get; set;} 
    public TileType type {get; set;} 

    protected Tile(int posX, int posY, TileType type){
        this.posX = posX;
        this.posY = posY;
        this.type = type;
    }
}

class GreenTile : Tile
{
    public GreenTile(int posX, int posY) : base(posX, posY, TileType.Green) {

    }
}

class VioletTile : Tile
{
    public VioletTile(int posX, int posY) : base(posX, posY, TileType.Violet) {

    }
}

class Chunck : IEquatable<Chunck>
{
    public static int size = 32;
    public Vector2 currPos {get; set;}
    public Vector2 realPos {get; private set;}
    public List<Tile> floorTiles {get; set;}
    public bool isVisible = false;

    public bool isRendered = false;
    public List<GameObject> renderedTiles {get; set;}

    public Chunck(List<Tile> floorTiles, Vector2 realPos) {
        this.floorTiles = floorTiles;
        this.realPos = realPos;
    }

    public bool Equals(Chunck other) {
        if (other == null) {
            return false;
        }
        return this.realPos.x == other.realPos.x && this.realPos.y == other.realPos.y;
    }


// does not support grid changes!!!!!
// only works with girdHigh of 9
    public void setCurrPos(float sizePerUnit) {
        Vector3 pos = renderedTiles[0].transform.localPosition;
        currPos = new Vector2(pos.x + size/2 - sizePerUnit/2, pos.y - size/2 + sizePerUnit/2);
    }

    public void translate(Vector3 dir, float sizePerUnit) {
        if(isRendered) {
            foreach(GameObject reneredTile in renderedTiles) {
                reneredTile.transform.Translate(dir);
            }
        }
        setCurrPos(sizePerUnit);
    }


    public void CheckIfVisible(Vector2 cameraPoint1, Vector2 cameraPoint2, float sizePerUnit) {

        float offset = sizePerUnit * (size/2);
        Vector2 chunkPoint1= new Vector2(currPos.x - offset, currPos.y + offset);
        Vector2 chunkPoint2= new Vector2(currPos.x + offset, currPos.y - offset);
        
        isVisible = Helper.doOverlap(
            cameraPoint1, cameraPoint2,
            chunkPoint1, chunkPoint2
        );
    }

    public List<Chunck> getNextChuncks() {
        List<Chunck> nextChuncks = new List<Chunck>();
        nextChuncks.AddRange( new List<Chunck> {
                Test.genChunck(new Vector2(currPos.x, currPos.y+1)),
                Test.genChunck(new Vector2(currPos.x+1, currPos.y)),
                Test.genChunck(new Vector2(currPos.x-1, currPos.y)),
                Test.genChunck(new Vector2(currPos.x, currPos.y-1))
            }
        );
        return nextChuncks;
    }
}


public class TileType
{
    private TileType(string value) { Value = value; }

    public string Value { get; private set; }

    public static TileType Green   { get { return new TileType("Green"); } }
    public static TileType Violet   { get { return new TileType("Violet"); } }
    public static TileType Error   { get { return new TileType("Error"); } }
}

class Test
{
    public static List<Tile> genRandTiles(int nrOfTiles) {
        List<Tile> genTiles = new List<Tile>();
        for(int i = 0; i < nrOfTiles; i++) {
            for(int j = 0; j < nrOfTiles; j++) {
                int rand = UnityEngine.Random.Range(0, 2);
                if(rand == 0) {
                    genTiles.Add(new GreenTile(i, j));
                    continue;
                }
                else if(rand == 1) {
                    genTiles.Add(new VioletTile(i, j));
                }
            }
        }
        return genTiles;
    } 

    public static Chunck genChunck(Vector2 pos) {
        return new Chunck(genRandTiles(Chunck.size), pos);
    }
}


class Helper
{
    public static bool doOverlap(Vector2 l1, Vector2 r1,
                          Vector2 l2, Vector2 r2)
    {
        // If one rectangle is on left side of other
        if (l1.x >= r2.x || l2.x >= r1.x)
        {
            return false;
        }
 
        // If one rectangle is above other
        if (r1.y >= l2.y || r2.y >= l1.y)
        {
            return false;
        }
        return true;
    }
}






