/**
 * 
 */
package codetest.tingtingtang.services;

import java.util.Random;

import codetest.tingtingtang.entities.ShipInfo;
import codetest.tingtingtang.enums.Direction;
import codetest.tingtingtang.enums.ShipType;

/**
 * @author TingTing Tang
 *
 * This factory pattern is used for generate the board and different types of ships.
 * 
 */
public class BattleShipObjectFactory 
{
	/*
	 * Generate battleship and destroyer with the random generated direction.
	 * 
	 */
	public static ShipInfo generateShip(ShipType shipType)
	{
		ShipInfo shipInfo = null;

		Random generator = new Random();
		Direction direction = null;
		
		switch (shipType)
		{
		case BATTLESHIP:
			int randB = generator.nextInt(2);
			if(randB == Direction.HORIZONTAL.getValue())
				direction = Direction.HORIZONTAL;
			if(randB == Direction.VERTICAL.getValue())
				direction = Direction.VERTICAL;
			shipInfo = new ShipInfo(5, direction);
			break;
			
		case DESTROYER:
			int randD = generator.nextInt(2);
			if(randD == Direction.HORIZONTAL.getValue())
				direction = Direction.HORIZONTAL;
			if(randD == Direction.VERTICAL.getValue())
				direction = Direction.VERTICAL;
				shipInfo = new ShipInfo(4, direction);
			break;
		}
		
		return shipInfo;
	}
	
	/*
	 * Generate the board with certain length and width.
	 */
	public static boolean[][] generateBoard(int length, int width)
	{
		return new boolean[length][width];
	}
}
