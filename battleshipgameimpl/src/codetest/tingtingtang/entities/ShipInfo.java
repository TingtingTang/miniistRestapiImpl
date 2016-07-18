/**
 * 
 */
package codetest.tingtingtang.entities;

import codetest.tingtingtang.enums.Direction;

/**
 * @author TingTing Tang
 * 
 * This entity corresponds to ShipInfo which include the length of the ship, the direction of the ship. And the starting point X and Y is included as placing the ship.
 * Meanwhile, the hitCounter attribute is to counter the number of hit on the specific ship in order to judge the ship will be sink or not.
 *
 */
public class ShipInfo 
{
	private int shipLength;
	private Direction shipDirection;
	private int startingPointX;
	private int startingPointY;
	private int hitCounter;
	
	public ShipInfo(int shipLength, Direction shipDirection)
	{
		this.shipLength = shipLength;
		this.shipDirection = shipDirection;
		this.hitCounter = 0;
	}

	/**
	 * @return the shipLength
	 */
	public int getShipLength()
	{
		return shipLength;
	}

	/**
	 * @param shipLength the shipLength to set
	 */
	public void setShipLength(int shipLength)
	{
		this.shipLength = shipLength;
	}

	/**
	 * @return the shipDirection
	 */
	public Direction getShipDirection() 
	{
		return shipDirection;
	}

	/**
	 * @param shipDirection the shipDirection to set
	 */
	public void setShipDirection(Direction shipDirection) 
	{
		this.shipDirection = shipDirection;
	}

	/**
	 * @return the startingPointX
	 */
	public int getStartingPointX()
	{
		return startingPointX;
	}

	/**
	 * @param startingPointX the startingPointX to set
	 */
	public void setStartingPointX(int startingPointX)
	{
		this.startingPointX = startingPointX;
	}

	/**
	 * @return the startingPointY
	 */
	public int getStartingPointY() 
	{
		return startingPointY;
	}

	/**
	 * @param startingPointY the startingPointY to set
	 */
	public void setStartingPointY(int startingPointY) 
	{
		this.startingPointY = startingPointY;
	}
	
	/**
	 * @return the startingPointY
	 */
	public int getHitCounter() 
	{
		return hitCounter;
	}

	/**
	 * @param startingPointY the startingPointY to set
	 */
	public void setHitCounter(int hitCounter) 
	{
		this.hitCounter = hitCounter;
	}

}
