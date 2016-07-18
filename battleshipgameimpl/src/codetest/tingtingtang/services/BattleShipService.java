/**
 * 
 */
package codetest.tingtingtang.services;

import java.util.List;
import java.util.Random;

import codetest.tingtingtang.entities.ShipInfo;
import codetest.tingtingtang.enums.Direction;
import codetest.tingtingtang.enums.RowOfBoard;
import codetest.tingtingtang.enums.TargetOutcome;

/**
 * @author TingTing Tang
 *
 *	This class is the service for placing the ships on the board 
 *	and analysing the user input on the board corresponding to the placed ships' location.
 *
 */
public class BattleShipService
{
	List<ShipInfo> ships = null;
	boolean[][] board = null;

	public BattleShipService(List<ShipInfo> ships, boolean[][] board)
	{
		this.ships = ships;
		this.board = board;
		placeShip();
	}
	
	/*
	 * This method is used for placing ships. The logic used as explained as follow.
	 * 1. Check ship direction, if horizontal
	 * 2. Set limit for random number generation for the starting point of the ships' X coordinate.
	 * 3. If vertical, do Step 2 for Y coordinate.
	 * 4. Step 1, 2, 3 are for the first ship.
	 * 5. For remaining ships, do Step 1 and 2.
	 * 6. For horizontal placement, check the generated X coordinate does not overlapped with existing ships.
	 * 7. If true, place the ship.
	 * 8. For vertical placement, performance Step 5 and 6, but by focusing on the Y coordinate.
	 * 
	 */
	private void placeShip()
	{
		Random randGenerator = new Random();
		int counter = 0;
		int startingPointX = 0;
		int startingPointY = 0;
		for(ShipInfo ship : ships)
		{
			if(counter == 0)
			{
				if(ship.getShipDirection() == Direction.HORIZONTAL)
				{
					startingPointX = randGenerator.nextInt(10);
					startingPointY = randGenerator.nextInt(10-ship.getShipLength()+1);
					ship.setStartingPointX(startingPointX);
					ship.setStartingPointY(startingPointY);
					for(int i = startingPointY; i < startingPointY+ship.getShipLength(); i++)
					{
						board[startingPointX][i] = true;
					}
				}
				else if(ship.getShipDirection() == Direction.VERTICAL)
				{
					startingPointX = randGenerator.nextInt(10-ship.getShipLength()+1);
					startingPointY = randGenerator.nextInt(10);
					ship.setStartingPointX(startingPointX);
					ship.setStartingPointY(startingPointY);
					for(int i = startingPointX; i < startingPointX+ship.getShipLength(); i++)
					{
						board[i][startingPointY] = true;
					}
				}
				counter++;
			}
			else
			{
				if(ship.getShipDirection() == Direction.HORIZONTAL)
				{
					startingPointX = randGenerator.nextInt(10);
					startingPointY = randGenerator.nextInt(10-ship.getShipLength()+1);
					int checkCounter = 0;
					while(checkCounter < ship.getShipLength())
					{
						if(board[startingPointX][startingPointY+checkCounter] == false)
						{
							checkCounter++;
						}
						else
						{
							startingPointX = randGenerator.nextInt(10);
							startingPointY = randGenerator.nextInt(10-ship.getShipLength()+1);
							checkCounter = 0;
						}
					}
					ship.setStartingPointX(startingPointX);
					ship.setStartingPointY(startingPointY);
					for(int j = startingPointY; j < startingPointY+ship.getShipLength(); j++)
					{
						board[startingPointX][j] = true;
					}
				}
				else if(ship.getShipDirection() == Direction.VERTICAL)
				{
					startingPointX = randGenerator.nextInt(10-ship.getShipLength()+1);
					startingPointY = randGenerator.nextInt(10);
					int checkCounter = 0;
					while(checkCounter < ship.getShipLength())
					{
						if(board[startingPointX+checkCounter][startingPointY] == false)
						{
							checkCounter++;
						}
						else
						{
							startingPointX = randGenerator.nextInt(10-ship.getShipLength()+1);
							startingPointY = randGenerator.nextInt(10);
							checkCounter = 0;
						}
					}
					ship.setStartingPointX(startingPointX);
					ship.setStartingPointY(startingPointY);
					for(int k = startingPointX; k < startingPointX+ship.getShipLength(); k++)
					{
						board[k][startingPointY] = true;
					}
				}
			}
		}
	}

	private int getRowValue(String valueInput)
	{
		int enumValue = 0;
		if(RowOfBoard.A.name().equalsIgnoreCase(valueInput))
		{
			enumValue = RowOfBoard.A.getValue();
		}

		if(RowOfBoard.B.name().equalsIgnoreCase(valueInput))
		{
			enumValue = RowOfBoard.B.getValue();
		}

		if(RowOfBoard.C.name().equalsIgnoreCase(valueInput))
		{
			enumValue = RowOfBoard.C.getValue();
		}

		if(RowOfBoard.D.name().equalsIgnoreCase(valueInput))
		{
			enumValue = RowOfBoard.D.getValue();
		}

		if(RowOfBoard.E.name().equalsIgnoreCase(valueInput))
		{
			enumValue = RowOfBoard.E.getValue();
		}

		if(RowOfBoard.F.name().equalsIgnoreCase(valueInput))
		{
			enumValue = RowOfBoard.F.getValue();
		}

		if(RowOfBoard.G.name().equalsIgnoreCase(valueInput))
		{
			enumValue = RowOfBoard.G.getValue();
		}

		if(RowOfBoard.H.name().equalsIgnoreCase(valueInput))
		{
			enumValue = RowOfBoard.H.getValue();
		}

		if(RowOfBoard.I.name().equalsIgnoreCase(valueInput))
		{
			enumValue = RowOfBoard.I.getValue();
		}

		if(RowOfBoard.J.name().equalsIgnoreCase(valueInput))
		{
			enumValue = RowOfBoard.J.getValue();
		}

		return enumValue;

	}
	
	
	public String analyseTarget(String userInput)
	{
		char[] location = userInput.toCharArray();
		String targetOutcome = null;

		int locationRow = getRowValue(Character.toString(location[0]));
		int locationCol = Character.getNumericValue(location[1]) - 1;

		if(ships.size() > 0)
		{
			if(board[locationRow][locationCol] == false)
			{
				targetOutcome = TargetOutcome.MISS.name();
			}
			else
			{
				for(ShipInfo ship : ships)
				{
					if(ship.getShipDirection() == Direction.HORIZONTAL)
					{
						if(ship.getStartingPointX() == locationRow)
						{
							if(ship.getStartingPointY() <= locationCol && ship.getStartingPointY()+ship.getShipLength() > locationCol)
							{
								int hitCounter = ship.getHitCounter() + 1;
								ship.setHitCounter(hitCounter);
								targetOutcome = TargetOutcome.HIT.name();
							}
						}
					}
					else if(ship.getShipDirection() == Direction.VERTICAL)
					{
						if(ship.getStartingPointY() == locationCol)
						{
							if(ship.getStartingPointX() <= locationRow && ship.getStartingPointX()+ship.getShipLength() > locationRow)
							{
								int hitCounter = ship.getHitCounter() + 1;
								ship.setHitCounter(hitCounter);
								targetOutcome = TargetOutcome.HIT.name();
							}
						}
					}
					if(ship.getHitCounter() == ship.getShipLength())
					{
						targetOutcome = TargetOutcome.SINK.name();
						ships.remove(ship);
						break;
					}
				}
			}
		}
		else
		{
			targetOutcome = TargetOutcome.WIN.name();
		}
		return targetOutcome;
	}
}
