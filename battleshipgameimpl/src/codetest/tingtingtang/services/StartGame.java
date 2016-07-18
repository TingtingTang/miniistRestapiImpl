/**
 * 
 */
package codetest.tingtingtang.services;

import java.util.ArrayList;
import java.util.List;
import java.util.Scanner;

import codetest.tingtingtang.entities.ShipInfo;
import codetest.tingtingtang.enums.RowOfBoard;
import codetest.tingtingtang.enums.ShipType;
import codetest.tingtingtang.enums.TargetOutcome;

/**
 * @author TingTing Tang
 *
 */
public class StartGame
{

	/**
	 * Initialise ships, board, service class and validate user's input.
	 */
	
	public static void main(String[] args) 
	{
		// TODO Auto-generated method stub
		List<ShipInfo> ships = new ArrayList<ShipInfo>();
		ShipInfo battleShip = BattleShipObjectFactory.generateShip(ShipType.BATTLESHIP);
		ShipInfo firstDestroyer = BattleShipObjectFactory.generateShip(ShipType.DESTROYER);
		ShipInfo secondDestroyer = BattleShipObjectFactory.generateShip(ShipType.DESTROYER);
		ships.add(battleShip);
		ships.add(firstDestroyer);
		ships.add(secondDestroyer);

		boolean[][] board = BattleShipObjectFactory.generateBoard(10, 10);
		BattleShipService battleShipService = new BattleShipService(ships, board);
		ArrayList<String> inputs = new ArrayList<String>();
		Scanner scanner = new Scanner(System.in);
		boolean flag = false;
		while(flag == false)
		{
			System.out.print("Enter input:");
			String userInput = scanner.next();
			if(validateInput(userInput, inputs) != null)
			{
				String outcome = battleShipService.analyseTarget(userInput);
				if(outcome.equalsIgnoreCase(TargetOutcome.WIN.name()))
					flag = true;
				
				System.out.println(outcome + "\n");
				inputs.add(userInput);
			}
			else
				System.out.println("Error in input. Please try again. \n\n");	
		}
		scanner.close();
	}

	private static String validateInput(String userInput, ArrayList<String> inputs)
	{
		//Check input is A - J
		char[] charInputs = userInput.toCharArray();

		if(charInputs.length > 3 && charInputs.length <= 1)
			return null;
		
		if(inputs.contains(userInput))
			return null;

		boolean characterValidated = validateRow(String.valueOf(charInputs[0]));
		boolean numeralValidated = false;
		if(charInputs.length == 2)
		{
			int number = Integer.parseInt(String.valueOf(charInputs[1]));
			if(number > 0 && number < 10)
				numeralValidated = true;
		}
		else if(charInputs.length == 3)
		{
			int number = Integer.parseInt(String.valueOf(charInputs[1])+String.valueOf(charInputs[2]));
			if(number == 10)
				numeralValidated = true;
		}

		if(characterValidated && numeralValidated)
			return userInput;
		else
			return null;

	}

	private static boolean validateRow(String value)
	{
		if(RowOfBoard.A.name().equalsIgnoreCase(value))
		{
			return true;
		}

		if(RowOfBoard.B.name().equalsIgnoreCase(value))
		{
			return true;
		}

		if(RowOfBoard.C.name().equalsIgnoreCase(value))
		{
			return true;
		}

		if(RowOfBoard.D.name().equalsIgnoreCase(value))
		{
			return true;
		}

		if(RowOfBoard.E.name().equalsIgnoreCase(value))
		{
			return true;
		}

		if(RowOfBoard.F.name().equalsIgnoreCase(value))
		{
			return true;
		}

		if(RowOfBoard.G.name().equalsIgnoreCase(value))
		{
			return true;
		}

		if(RowOfBoard.H.name().equalsIgnoreCase(value))
		{
			return true;
		}

		if(RowOfBoard.I.name().equalsIgnoreCase(value))
		{
			return true;
		}

		if(RowOfBoard.J.name().equalsIgnoreCase(value))
		{
			return true;
		}

		return false;
	}

}
